package sqlite

import (
	"backend/article"
	"backend/log"
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
	"github.com/gofrs/uuid"
)

var l = log.Logger.With().Str("component", "sqlite").Logger()

type DbManager struct {
	db   *sql.DB
	path string
}

func Init(filepath string) *DbManager {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return nil
	}

	// Create a table
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS articles (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			url TEXT NOT NULL,
			tags TEXT,
			friendly_url TEXT,
			creation_timestamp INTEGER,
			edit_timestamp INTEGER,
			meta_description TEXT,
			published INTEGER,
			integrity_hash TEXT
		);
	`
	_, err = db.Exec(createTableSQL)

	if err != nil {
		l.Err(err)
		return nil
	}

	l.Debug().Msg("Database table has been created.")

	return &DbManager{
		db:   db,
		path: filepath,
	}
}

func (db *DbManager) GetArticlePath(uuid string) (string, error) {

	row := db.db.QueryRow("SELECT url FROM articles WHERE id == ? LIMIT 1", uuid)
	l.Debug().Msg("Querying database for article path based on id property. [GetArticlePath]")
	var url string
	err := row.Scan(&url)
	l.Debug().Msg("Parsing the results of article path query. [GetArticlePath]")
	if err != nil {
		return "", err
	}

	return url, nil
}

// adjusted to the new database
func (db *DbManager) GetArticles() ([]*article.Article, error) {
	var articles []*article.Article

	rows, err := db.db.Query("SELECT * FROM articles")
	for rows.Next() {

		var title, url, uuid, friendlyUrl, tags string
		var creationTimestamp uint64
		var editTimestamp *uint64
		var metaDescription *string
		var published bool
		if err := rows.Scan(&uuid, &title, &url, &tags, &friendlyUrl, &creationTimestamp, &editTimestamp, &metaDescription); err != nil {
			return nil, err
		}

		a := article.NewFromScratch(uuid, title, url, article.CsvToTags(tags), friendlyUrl, creationTimestamp, editTimestamp, metaDescription, published)
		if a != nil {
			l.Warn().Msg("Failed to create new article. Most likely due to UUID error. The article will be omitted")
			continue
		}

		articles = append(articles, a)
		if err != nil {
			return nil, err
		}
		l.Debug().Msg(fmt.Sprintf("Scanned the article - %s [GetArticles]", title))
	}
	return articles, nil
}

func (db *DbManager) IfUrlExist(a article.Article) bool {
	row := db.db.QueryRow("SELECT * FROM articles WHERE url == ?", a.Url)
	l.Debug().Msg("Querying database to check if url exists in any row. [IfUrlExist]")
	var box string
	var dst article.Article
	err := row.Scan(&dst.Uuid, &dst.Title, &dst.Url, &box)
	if err != nil {
		l.Debug().Msg(fmt.Sprintf("Error raised during scanning rows: %s. [IfUrlExist]", err.Error()))
		return false
	}
	return true
}

// Handles user input
func (db *DbManager) IfRowExist(a article.Article) bool {
	row := db.db.QueryRow("SELECT * FROM articles WHERE id == ?", a.Uuid)
	l.Debug().Msg("Querying database to check if row exists. [IfRowExist]")
	var dst article.Article
	err := row.Scan(&dst.Uuid, &dst.Title, &dst.Url, &dst.Tags)
	if err != nil {
		l.Debug().Msg(fmt.Sprintf("Error raised during scanning rows: %s. [IfRowExist]", err.Error()))
		return false
	} else if dst.Uuid == a.Uuid && dst.Url == a.Url {
		l.Debug().Msg("The row exist and the ID and path matches the one in the database. [IfRowExist]")
		return true
	}
	l.Debug().Msg(fmt.Sprintf("The suspicious manipulation of protobuf. UUID: %s, should be: %s. Path: %s, should be: %s. [IfRowExist]", a.Uuid, dst.Uuid, a.Url, dst.Url))
	return false
}

func (db *DbManager) UpdateRecord(title string, url string, tags string) error {

	row := db.db.QueryRow("SELECT id FROM articles WHERE url == ?", url)
	l.Debug().Msg("Querying database to check if url exists in any row. [IfUrlExist]")
	var uuid string
	err := row.Scan(&uuid)
	if err != nil {
		l.Debug().Msg(fmt.Sprintf("Error raised during scanning rows: %s. [IfRowExist]", err.Error()))
		return err
	}
	res, err := db.db.Exec(`
	UPDATE articles
	SET title = ?,
		url = ?,
		tags = ?
	WHERE
		id = ?`,
		title,
		url,
		tags,
		uuid,
	)
	if err != nil {
		return err
	}
	l.Debug().Msg(fmt.Sprintf("Record was created. ID: %s, Title: %s, Url: %s, Tags: %s", uuid, title, url, tags))
	//The RowsAffected always return nil for error.
	nAffected, _ := res.RowsAffected()
	if nAffected > 1 {
		l.Warn().Msg(fmt.Sprintf("More than 1 record was created. ID: %s, Title: %s, Url: %s, Tags: %s", uuid, title, url, tags))
	}

	return nil
}

func (db *DbManager) CreateRecord(title string, url string, tags string) error {
	id, err := uuid.NewV6()
	if err != nil {
		return err
	}
	result, err := db.db.Exec(`INSERT INTO articles (id, title, url, tags)
									VALUES
									(?, ?, ? ,?)`, id, title, url, tags)
	if err != nil {
		return err
	}
	l.Debug().Msg(fmt.Sprintf("Record was created. ID: %s, Title: %s, Url: %s, Tags: %s", id, title, url, tags))
	//The RowsAffected always return nil for error.
	res, _ := result.RowsAffected()
	if res > 1 {
		l.Warn().Msg(fmt.Sprintf("More than 1 record was created. ID: %s, Title: %s, Url: %s, Tags: %s", id, title, url, tags))
	}
	return nil
}
