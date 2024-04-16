package article

import (
	"backend/log"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

var l = log.Logger.With().Str("component", "article").Logger()

type Article struct {
	Uuid              string
	Title             string
	Url               string
	Tags              []string
	FriendlyUrl       string
	CreationTimestamp uint64
	EditTimestamp     *uint64
	MetaDescription   *string
	Published         bool
	IntegrityHash     string
}

func New(title string, url string, tags []string, friendlyUrl string, metaDescription string, published bool) *Article {
	id, err := uuid.NewV6()
	if err != nil {
		l.Error().Msg(fmt.Sprintf("Error creating new article [NewArticle]: %s", err))
		return nil
	}

	var editTimestamp uint64 = uint64(time.Now().Unix())

	return &Article{
		Uuid:              id.String(),
		Title:             title,
		Url:               url,
		Tags:              tags,
		FriendlyUrl:       friendlyUrl,
		CreationTimestamp: uint64(time.Now().Unix()),
		EditTimestamp:     &editTimestamp,
		MetaDescription:   &metaDescription,
		Published:         published,
	}
}

func NewFromScratch(uuid string, title string, url string, tags []string, friendlyUrl string, creationTimestamp uint64, editTimestamp *uint64, metaDescription *string, published bool) *Article {
	tmp := &Article{
		Uuid:              uuid,
		Title:             title,
		Url:               url,
		Tags:              tags,
		FriendlyUrl:       friendlyUrl,
		CreationTimestamp: creationTimestamp,
		EditTimestamp:     editTimestamp,
		MetaDescription:   metaDescription,
		Published:         published,
	}

	tmp.CalculateArticleIntegrityHash()

	return tmp
}

func (a Article) CalculateArticleIntegrityHash() {
	editTime := ""
	metaDescription := ""
	if a.EditTimestamp != nil {
		editTime = fmt.Sprint(*a.EditTimestamp)
	}

	if a.MetaDescription != nil {
		metaDescription = *a.MetaDescription
	}

	data := a.Uuid + a.Title + a.Url + a.FriendlyUrl + fmt.Sprint(a.CreationTimestamp) + string(editTime) + string(metaDescription)

	a.IntegrityHash = hex.EncodeToString(sha256.Sum256([]byte(data)))
}

func (a Article) DebugPrint() {
	fmt.Printf("%s, %s, %s", a.Uuid, a.Title, a.Url)

}

func LegacyNew(uuid string, title string, url string, tags []string) Article {
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
