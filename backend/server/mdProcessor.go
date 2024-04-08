package server

import (
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Extension int

const (
	HTML Extension = iota
	MD
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.BackslashLineBreak
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func Read(file string) ([]byte, error) {
	mdFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer mdFile.Close()

	mdFileStats, err := mdFile.Stat()
	if err != nil {
		return nil, err
	}

	buff := make([]byte, mdFileStats.Size())

	_, err = mdFile.Read(buff)
	if err != nil {
		return nil, err
	}

	return buff, nil

}

func Write(file string, content []byte) error {
	htmlFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer htmlFile.Close()

	_, err = htmlFile.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func MdFileToHTMLFile(src string, dst string) error {

	srcContent, err := Read(src)
	if err != nil {
		return err
	}

	dstContent := mdToHTML(srcContent)

	err = Write(dst, dstContent)
	if err != nil {
		return err
	}

	return nil
}
