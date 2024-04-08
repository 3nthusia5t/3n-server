package sqlite

import (
	"fmt"
	"testing"

	_ "github.com/glebarez/go-sqlite"
)

func TestQueryRow(t *testing.T) {

	localDb := Init("../database.db")

	url, err := localDb.GetArticlePath("1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(url)
}
