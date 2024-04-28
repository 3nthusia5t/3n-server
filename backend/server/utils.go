package server

import (
	"backend/article"
	"backend/model"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/proto"
)

func SerializeArticles(articles []*article.Article) ([]byte, error) {
	var serializedMessage []byte
	var protoArticles []*model.Article
	for _, elem := range articles {

		if !elem.Published {
			continue
		}

		tempEditTimestamp := int64(*elem.EditTimestamp)
		article := &model.Article{
			Uuid:              elem.Uuid,
			Title:             elem.Title,
			Url:               elem.Url,
			Tags:              elem.TagsToCsv(),
			FriendlyUrl:       elem.FriendlyUrl,
			CreationTimestamp: int64(elem.CreationTimestamp),
			EditTimestamp:     &tempEditTimestamp,
			MetaDescription:   elem.MetaDescription,
			Published:         elem.Published,
		}
		protoArticles = append(protoArticles, article)
	}

	pbArticleList := &model.Articles{
		ListOfArticles: protoArticles,
	}
	msgBuff, err := proto.Marshal(pbArticleList)
	if err != nil {
		return nil, err
	}
	serializedMessage = msgBuff
	return serializedMessage, nil
}

func UnserializeArticle(serializedMessage []byte) (article.Article, error) {
	m := model.Article{}
	err := proto.Unmarshal(serializedMessage, &m)
	if err != nil {
		return article.Article{}, err
	}

	article := article.Article{
		Uuid:  m.Uuid,
		Title: m.Title,
		Url:   m.Url,
		Tags:  article.CsvToTags(m.Tags),
	}

	return article, nil

}

func findAndCopyImages(sourceDir, destinationDir string) error {
	imageExtensions := map[string]bool{".svg": true, ".jpg": true, ".jpeg": true, ".png": true}

	err := filepath.Walk(sourceDir, func(srcPath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//Check if path is part of other path. Don't copy dst dir to dst dir.
		relPath, err := filepath.Rel(srcPath, destinationDir)

		if err != nil {
			return nil
		}

		if relPath == ".." || relPath == "." {
			return nil
		}

		if fileInfo.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(fileInfo.Name()))
		if imageExtensions[ext] {
			destPath := filepath.Join(destinationDir, fileInfo.Name())
			if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
				return err
			}

			err := copyFile(srcPath, destPath)
			if err != nil {
				fmt.Printf("Error copying file %s: %s\n", srcPath, err)
			} else {
				fmt.Printf("Copied file: %s\n", srcPath)
			}
		}

		return nil
	})

	return err
}

func TranscompileArticles(srcPath string, dstPath string) error {
	l.Debug().Msg(fmt.Sprintf("Copying the source directory (%s) to the destination directory (%s). [TranscompileArticles]", srcPath, dstPath))
	err := copyDir(srcPath, dstPath)
	if err != nil {
		return err
	}

	var toDelete []string
	err = filepath.Walk(dstPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileExtension := filepath.Ext(path)

		if fileExtension == ".md" {

			l.Debug().Msg(fmt.Sprintf("Walking through the destination directory. Current working directory: %s. [TranscompileArticles]", path))
			content, err := Read(path)
			if err != nil {
				return err
			}
			htmlContent := mdToHTML(content)
			l.Debug().Msg("Transcompilating the markdown into HTML completed. [TranscompileArticles]")
			outputFile := changeFileExtension(path, "html")
			l.Debug().Msg(fmt.Sprintf("HTML file - %s will be created. [TranscompileArticles]", outputFile))
			err = Write(outputFile, htmlContent)
			if err != nil {
				l.Err(err)
				return err
			}
			toDelete = append(toDelete, path)
		}

		return nil
	})
	for _, file := range toDelete {
		err := os.Remove(file)
		if err != nil {
			l.Debug().Msg(err.Error())
		}
	}

	if err != nil {
		return err
	}
	return nil
}

func EnumerateArticles(path string) ([]*article.Article, error) {
	whitelist := []string{
		"README.md",
		".git",
		"LICENSE",
	}

	type meta struct {
		Title             string   `json:"title"`
		Tags              []string `json:"tags"`
		FriendlyUrl       string   `json:"friendly_url"`
		CreationTimestamp uint64   `json:"timestamp"`
		MetaDescription   string   `json:"meta_description"`
		Published         bool     `json:"published"`
		Path              string
	}

	var al []*article.Article
	var ml []meta

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

		//go through all kinds of filters
		if err != nil {
			return err
		}

		for i := range whitelist {
			if strings.Contains(path, whitelist[i]) {
				return nil
			}
		}

		if info.IsDir() {
			return nil
		}

		dirName, fileName := filepath.Split(path)

		//parse metadata files
		if fileName == "metadata" {
			content, err := Read(path)
			if err != nil {
				l.Error().Msg(fmt.Sprintf("Error reading file [EnumerateArticles]: %s", err.Error()))
				return err
			}
			l.Debug().Msg(fmt.Sprintf("Successfully read file [EnumerateArticles]: %s", path))

			var res meta
			err = json.Unmarshal(content, &res)
			if err != nil {
				l.Error().Msg(fmt.Sprintf("Error unmarshaling JSON content [EnumerateArticles]: %s", err.Error()))
				return err
			}
			l.Debug().Msg(fmt.Sprintf("Successfully unmarshal JSON content [EnumerateArticles]: %v", res))

			res.Path, err = findFileByExtensionInDir(dirName, "html")
			if err != nil {
				return err
			}
			ml = append(ml, res)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	//translate meta into article
	for _, metaItem := range ml {
		a := article.New(metaItem.Title, metaItem.Path, metaItem.Tags, metaItem.FriendlyUrl, metaItem.CreationTimestamp, metaItem.MetaDescription, metaItem.Published)
		al = append(al, a)
	}

	return al, nil
}

func findFileByExtensionInDir(dir, extension string) (string, error) {
	var file string
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(), "."+extension) {
			file = filepath.Join(dir, fileInfo.Name())
		}
	}

	return file, nil
}

func copyDir(src, dest string) error {
	// Get information about the source directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	err = os.MkdirAll(dest, srcInfo.Mode())
	if err != nil {
		return err
	}

	// Read the contents of the source directory
	dir, err := os.Open(src)
	if err != nil {
		return err
	}
	defer dir.Close()

	// Get the list of files and subdirectories in the source directory
	fileInfos, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	// Copy each file and subdirectory
	for _, fileInfo := range fileInfos {
		srcPath := filepath.Join(src, fileInfo.Name())
		destPath := filepath.Join(dest, fileInfo.Name())

		if fileInfo.IsDir() {
			// Recursively copy subdirectories
			err = copyDir(srcPath, destPath)
			if err != nil {
				fmt.Printf("Error copying directory %s: %s\n", srcPath, err)
			}
		} else {
			// Copy files
			err = copyFile(srcPath, destPath)
			if err != nil {
				fmt.Printf("Error copying file %s: %s\n", srcPath, err)
			}
		}
	}

	return nil
}

func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func changeFileExtension(filePath, newExtension string) string {
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	newFileName := base[:len(base)-len(filepath.Ext(base))] + "." + newExtension
	newPath := filepath.Join(dir, newFileName)
	return newPath
}
