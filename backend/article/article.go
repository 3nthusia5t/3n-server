package article

import (
	"fmt"
	"strings"

	uuid "github.com/google/uuid"
)

type Article struct {
	Uuid  string
	Title string
	Url   string
	Tags  []string
}

func NewArticle(title string, url string, tags []string) Article {
	return Article{
		Uuid:  uuid.NewString(),
		Title: title,
		Url:   url,
		Tags:  tags,
	}
}

func New(uuid string, title string, url string, tags []string) Article {
	return Article{
		Uuid:  uuid,
		Title: title,
		Url:   url,
		Tags:  tags,
	}
}

func CsvToTags(CsvTags string) []string {
	var tags []string
	splittedTags := strings.Split(CsvTags, ",")
	if len(splittedTags) == 0 {
		return tags
	}
	//Sanitize
	for _, tag := range splittedTags {
		if tag[0] == byte(' ') {
			tag = tag[1:]
		}
		if tag[len(tag)-1] == byte(' ') {
			tag = tag[:(len(tag) - 2)]
		}
		tags = append(tags, tag)
	}
	return tags
}

func TagsToCsv(tags []string) string {
	var CsvTags string
	for i, tag := range tags {
		if i == len(tags)-1 {
			CsvTags += tag
			continue
		}
		CsvTags += fmt.Sprintf("%s, ", tag)

	}
	return CsvTags
}

func (a Article) DebugPrint() {
	fmt.Printf("%s, %s, %s", a.Uuid, a.Title, a.Url)

}
