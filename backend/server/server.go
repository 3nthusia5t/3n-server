package server

import (
	"backend/article"
	"backend/log"
	"backend/sqlite"
	"fmt"
	"net/http"
	"time"
)

// GLOBAL VARS
var l = log.Logger.With().Str("component", "server").Logger()
var gArticles []*article.Article
var IsDev bool = false

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
func ServeApp(staticContentPath string, imagesContentPath string, externalContentPath string, databasePath string, tlsCertPath string, tlsKeyPath string, isDev *bool) {
	// it's useful to be able to restart server. This is a label for goto statements

	//make it global scope, so other packages can also use it.
	IsDev = *isDev

	//initialize the database
	DbManager := sqlite.Init(databasePath)
	l.Info().Msg(fmt.Sprintf("database path: %s", databasePath))
	if DbManager == nil {
		l.Fatal().Msg("Failed to initialize the database")

	} else {
		l.Info().Msg("Successfully initialized the database")
	}

	var err error
	gArticles, err = DbManager.GetArticles()
	go LoadArticlesToMemory(DbManager)
	if err != nil {
		l.Error().Msg(err.Error())
	}

	//Serving static website
	staticServer := http.FileServer(http.Dir(staticContentPath))
	imgServer := http.FileServer(http.Dir(imagesContentPath))
	http.Handle("/", http.StripPrefix("/", staticServer))
	http.Handle("/images/", http.StripPrefix("/images/", imgServer))

	//Handling API calls
	http.HandleFunc("/GetAllArticles", GetAllArticlesHandler)
	http.HandleFunc("/GetChosenArticle", GetChosenArticleHandler(DbManager))

	httpDone := make(chan struct{})

	go ServeHttp(httpDone, nil)

	<-httpDone
}

func ServeHttp(done chan struct{}, handler http.Handler) {
	defer close(done)
	err := http.ListenAndServe(":http", handler)
	if err != nil {
		l.Err(err).Msg("HTTP server error")
	}
	done <- struct{}{}
}

// ServeHttps starts and serves the HTTPS server
func ServeHttps(done chan struct{}, tlsCertPath string, tlsKeyPath string) {
	defer close(done)
	err := http.ListenAndServeTLS(":https", tlsCertPath, tlsKeyPath, nil)
	if err != nil {
		l.Err(err).Msg("HTTPS server error")
	}
	done <- struct{}{}
}

func UpdateApp(articleContentPath string, databasePath string) {
	//Enumarate articles and find new
	al, err := EnumerateArticles(articleContentPath)
	l.Info().Msg(fmt.Sprintf("%v", al))
	if err != nil {
		l.Error().Msg(err.Error())
	}
	l.Debug().Msg(fmt.Sprintf("Successfully enumerated articles [UpdateApp]: %v", al))

	DbManager := sqlite.Init(databasePath)
	for _, a := range al {
		fmt.Println(a.FriendlyUrl)
		/*
			We're checking URL instead of UUID, because the metadata files in 3n-articles doesn't contain uuid.
			Therefore UUID is completely random and doesn't relate to the one that is already stored in database.
		*/
		if DbManager.IfUrlExist(*a) {
			//EditTimestamp assignment inside the function
			DbManager.UpdateRecord(a)
		} else {
			DbManager.CreateRecord(*a)
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
