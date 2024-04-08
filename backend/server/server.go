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

func LoadArticlesToMemory(localDb *sqlite.LocalDB) {
	for true {
		time.Sleep(30 * time.Minute)
		var err error
		gArticles, err = localDb.GetArticles()
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
	localDb := sqlite.Init(databasePath)
	if localDb == nil {
		l.Fatal().Msg("Failed to initialize the database")
	} else {
		l.Info().Msg("Successfully initialized the database")
	}

	var err error
	gArticles, err = localDb.GetArticles()
	go LoadArticlesToMemory(localDb)
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
	http.HandleFunc("/GetChosenArticle", GetChosenArticleHandler(localDb))

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
	localDb := sqlite.Init(databasePath)

	for _, a := range al {
		print(localDb.IfUrlExist(a))
		if localDb.IfUrlExist(a) {
			localDb.UpdateRecord(a.Title, a.Url, article.TagsToCsv(a.Tags))
		} else {
			localDb.CreateRecord(a.Title, a.Url, article.TagsToCsv(a.Tags))
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

func Test(staticContentPath string, imagesContentPath string, externalContentPath string) {
	return
}
