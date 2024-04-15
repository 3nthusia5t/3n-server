package sqlite

import (
	"backend/sqlite"
	"fmt"
	"testing"

	_ "github.com/glebarez/go-sqlite"
)

func TestQueryRow(t *testing.T) {

	DbManager := sqlite.Init("../database.db")

	url, err := DbManager.GetArticlePath("1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(url)
}
