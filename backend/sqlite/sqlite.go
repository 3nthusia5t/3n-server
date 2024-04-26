package sqlite

import (
	"backend/article"
	"backend/log"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/glebarez/go-sqlite"
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

	rows, err := db.db.Query("SELECT id, title, url, tags, friendly_url, creation_timestamp, edit_timestamp, meta_description, published FROM articles")
	for rows.Next() {

		var title, url, uuid, friendlyUrl, tags string
		var creationTimestamp uint64
		var editTimestamp *uint64
		var metaDescription *string
		var published bool
		if err := rows.Scan(&uuid, &title, &url, &tags, &friendlyUrl, &creationTimestamp, &editTimestamp, &metaDescription, &published); err != nil {
			return nil, err
		}

		a := article.NewFromScratch(uuid, title, url, article.CsvToTags(tags), friendlyUrl, creationTimestamp, editTimestamp, metaDescription, published)
		if a == nil {
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
	row := db.db.QueryRow("SELECT id, integrity_hash FROM articles WHERE url == ?", a.Url)
	l.Debug().Msg("Querying database to check if url exists in any row. [IfUrlExist]")
	var dst article.Article
	err := row.Scan(&dst.Uuid, &dst.IntegrityHash)

	if a.IntegrityHash != dst.IntegrityHash {
		l.Debug().Msg(fmt.Sprintf("Provided integrity hash %s, does not much the original one: %s", a.IntegrityHash, dst.IntegrityHash))
		return false
	}

	if err != nil {
		l.Debug().Msg(fmt.Sprintf("Error raised during scanning rows: %s. [IfUrlExist]", err.Error()))
		return false
	}

	l.Debug().Msg("URL exists in the database.")
	return true
}

// Unused
// TODO:
// Make it work with new schema.
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

func (db *DbManager) UpdateRecord(a *article.Article) error {

	row := db.db.QueryRow("SELECT id FROM articles WHERE url = ?", a.Url)
	var uuid string
	err := row.Scan(&uuid)
	l.Debug().Msg(fmt.Sprintf("Updating the record with UUID %s [UpdateRecord]", uuid))

	if err != nil {
		l.Debug().Msg(fmt.Sprintf("Error raised during scanning rows: %s. [UpdateRecord]", err.Error()))
		return err
	}

	fmt.Println(a.FriendlyUrl)

	updatedEditTimestamp := uint64(time.Now().Unix())
	a.EditTimestamp = &updatedEditTimestamp
	fmt.Println(uuid)
	fmt.Println(a.Uuid)
	res, err := db.db.Exec(`
	UPDATE articles
	SET title = ?,
		url = ?,
		tags = ?,
		friendly_url = ?,
		creation_timestamp = ?,
		edit_timestamp = ?,
		meta_description = ?,
		published = ?,
		integrity_hash = ?
	WHERE
		id = ?`,
		a.Title,
		a.Url,
		a.TagsToCsv(),
		a.FriendlyUrl,
		a.CreationTimestamp,
		a.EditTimestamp,
		a.MetaDescription,
		a.Published,
		a.IntegrityHash,
		uuid,
	)
	if err != nil {
		l.Debug().Msg(err.Error())
		return err
	}
	l.Debug().Msg(fmt.Sprintf("Record was updated: %s", a.DebugPrint()))
	//The RowsAffected always return nil for error.
	nAffected, _ := res.RowsAffected()
	fmt.Println(nAffected)
	if nAffected > 1 {
		l.Warn().Msg(fmt.Sprintf("More than 1 record was created: %s", a.DebugPrint()))
	}

	return nil
}

// fix EnumerateArticles in utils.go before fixing this one.

func (db *DbManager) CreateRecord(a article.Article) error {

	result, err := db.db.Exec(`
	INSERT INTO articles
		(
		id,
		title,
		url,
		tags,
		friendly_url,
		creation_timestamp,
		edit_timestamp,
		meta_description,
		published,
		integrity_hash)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.Uuid,
		a.Title,
		a.Url,
		a.TagsToCsv(),
		a.FriendlyUrl,
		a.CreationTimestamp,
		a.EditTimestamp,
		a.MetaDescription,
		a.Published,
		a.IntegrityHash,
	)

	if err != nil {
		return err
	}
	l.Debug().Msg(fmt.Sprintf("Record was created: %s", a.DebugPrint()))

	//The RowsAffected always return nil for error.
	res, _ := result.RowsAffected()
	if res > 1 {
		l.Warn().Msg(fmt.Sprintf("More than 1 record was created: %s", a.DebugPrint()))
	}
	return nil
}
