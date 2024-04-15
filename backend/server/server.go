package server

import (
	"backend/article"
	"backend/log"
	"backend/sqlite"
	"fmt"
	"net/http"
	"os"
	"time"
)

var l = log.Logger.With().Str("component", "server").Logger()

// GLOBAL VARS
var gArticles []article.Article

func LoadArticlesToMemory(DbManager *sqlite.DbManager) {
	for {
		time.Sleep(30 * time.Minute)
		var err error
		gArticles, err = DbManager.GetArticles()
		if err != nil {
			l.Warn().Err(err)
			continue
		}

	}
}

// update the tls part
func ServeApp(staticContentPath string, imagesContentPath string, externalContentPath string, databasePath string, tlsCertPath string, tlsKeyPath string) {
RESTART: // it's useful to be able to restart server. This is a label for goto statements

	//check if cert exists
	if _, err := os.Stat(tlsCertPath); os.IsNotExist(err) {
		l.Warn().Msg("Cert not found, retrying in 5 minutes")
		time.Sleep(5 * time.Minute)
		goto RESTART
	}

	//check if key exists
	if _, err := os.Stat(tlsKeyPath); os.IsNotExist(err) {
		time.Sleep(5 * time.Minute)
		goto RESTART
	}

	//initialize the database
	DbManager := sqlite.Init(databasePath)
	if DbManager == nil {
		l.Fatal().Msg("Failed to initialize the database")

	} else {
		l.Info().Msg("Successfully initialized the database")
	}

	var err error
	gArticles, err = DbManager.GetArticles()
	go LoadArticlesToMemory(DbManager)
	if err != nil {
		l.Err(err)
	}

	//Serving static website
	staticServer := http.FileServer(http.Dir(staticContentPath))
	imgServer := http.FileServer(http.Dir(imagesContentPath))
	http.Handle("/", http.StripPrefix("/", staticServer))
	http.Handle("/images/", http.StripPrefix("/images/", imgServer))

	//Handling API calls
	http.HandleFunc("/GetAllArticles", GetAllArticlesHandler)
	http.HandleFunc("/GetChosenArticle", GetChosenArticleHandler(DbManager))

	//Starting the server
	l.Err(http.ListenAndServeTLS(":https", tlsCertPath, tlsKeyPath, nil))
}

func UpdateApp(articleContentPath string, databasePath string) {
	//Enumarate articles and find new
	al, err := EnumerateArticles(articleContentPath)
	l.Debug().Msg(fmt.Sprint(al))
	if err != nil {
		l.Error().Msg(err.Error())
	}
	DbManager := sqlite.Init(databasePath)

	for _, a := range al {
		print(DbManager.IfUrlExist(a))
		if DbManager.IfUrlExist(a) {
			DbManager.UpdateRecord(a.Title, a.Url, article.TagsToCsv(a.Tags))
		} else {
			DbManager.CreateRecord(a.Title, a.Url, article.TagsToCsv(a.Tags))
		}
	}

}

func TranscompileApp(srcContentPath string, dstContentPath string) {
	err := TranscompileArticles(srcContentPath, dstContentPath)
	if err != nil {
		l.Err(err)
	}
	err = findAndCopyImages(dstContentPath, dstContentPath+"/images")
	if err != nil {
		l.Err(err)
	}
}

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	// Get the host and requested URL path
	host := r.Host
	url := "https://" + host + r.URL.Path

	// Redirect to the HTTPS version
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func Test(staticContentPath string, imagesContentPath string, externalContentPath string) {
	return
}
