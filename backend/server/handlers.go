package server

import (
	"backend/sqlite"
	"fmt"
	"io"
	"net/http"
)

func GetChosenArticleHandler(db *sqlite.DbManager) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if IsDev {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			if r.URL.Path != "/GetChosenArticle" {
				http.Error(w, "403 not found", http.StatusNotFound)
			}

			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Methods", "POST")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			if r.Method != "POST" {
				http.Error(w, "Bad method", http.StatusBadRequest)
				l.Err(fmt.Errorf("the request had wrong method %s", r.Method))
				return
			}

		}

		//What happens when message will be bigger that 1024bytes?
		buff := make([]byte, 1024)
		n, err := r.Body.Read(buff)
		buff = buff[:n]

		if err != nil && err != io.EOF {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			l.Err(err)
			return
		}

		article, err := UnserializeArticle(buff)
		l.Debug().Msg(fmt.Sprintf("Successfully unserialized article %s", article.Title))

		if err != nil {
			l.Err(err)
			return
		}
		path, err := db.GetArticlePath(article.Uuid)
		if err != nil {
			l.Err(err)
			return
		}
		http.ServeFile(w, r, path)
		l.Info().Msg("Request has been handled successfully [GetChosenArticleHandler]")
	}
}

func GetAllArticlesHandler(w http.ResponseWriter, r *http.Request) {
	if IsDev {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			w.WriteHeader(http.StatusNoContent)
			l.Debug().Msg(fmt.Sprintf("Handled CORS %s request", r.Method))
			return
		}
	}

	msg, err := SerializeArticles(gArticles)
	if err != nil {
		l.Err(err)
		return
	}

	l.Debug().Msg(fmt.Sprintf("Serialized given articles successfuly: %v [GetAllArticlesHandler]", gArticles))

	w.Header().Set("Content-Type", "application/octet-stream")

	if r.Method != "GET" {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		l.Err(fmt.Errorf("request has wrong method"))
		return
	}

	bWritten, err := w.Write(msg)
	l.Info().Msg(fmt.Sprintf("Response with %d bytes of size has been sent. [GetAllArticlesHandler]", bWritten))
	if err != nil {
		l.Err(err)
		return
	}
	l.Info().Msg("Request has been handled successfully")
}

/*
func MdParserHandler(MdArticlesPath string, HTMLArticlesPath string) error {
	MdArticles, err := os.ReadDir(MdArticlesPath)
	fmt.Println(MdArticles)
	if err != nil {
		return err
	}

	for _, a := range MdArticles {
		if a.IsDir() == true {
			continue
		}
		HtmlPath := HTMLArticlesPath + ChangeFileExtension()
		MdPath := MdArticlesPath + a.Name()
		if ValidPath(MdPath) {
			MdFileToHTMLFile(MdPath)
		}
	}
	return nil
}
*/
