package cms

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

type CMSPageMeta struct {
	Title         string    `bson:"title" json:"title"`
	CreatorUserID string    `bson:"creator_user_id" json:"creator_user_id"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
}

type CMSPagePreferences struct {
	AuthorVisible    bool `bson:"author_visible" json:"author_visible"`
	UpdatersVisible  bool `bson:"updaters_visible" json:"updaters_visible"`
	Dynamic          bool `bson:"dynamic" json:"dynamic"` // Does the page include GoLang-Template content?
	URLCaseSensitive bool `bson:"url_case_sensitive" json:"url_case_sensitive"`
	UseMarkdown      bool `bson:"markdown" json:"markdown"`
}

type CMSPage struct {
	URL         string             `bson:"url" json:"url"`
	Meta        CMSPageMeta        `bson:"meta" json:"meta"`
	Updates     []CMSPageUpdate    `bson:"updates" json:"updates"`
	Preferences CMSPagePreferences `bson:"preferences" json:"preferences"`
	Content     string             `bson:"content" json:"content"`
}

func (p *CMSPage) CreateFilter() *bson.M {
	return &bson.M{
		"url": p.URL,
	}
}

var (
	mdCache = cache.New(60*time.Minute, 12*time.Minute)
)

func (p *CMSPage) GetContent() (content string) {
	content = p.Content

	if p.Preferences.UseMarkdown {
		log.Println("-> Uses markdown")
		// markdown
		if c, found := mdCache.Get(p.Content); found {
			content = c.(string)
			log.Println("-> from cache")
		} else {
			b := []byte(content)

			mdParser := parser.NewWithExtensions(
				parser.CommonExtensions | parser.AutoHeadingIDs |
					parser.Footnotes | parser.FencedCode |
					parser.DefinitionLists | parser.HardLineBreak,
			)
			html := markdown.ToHTML(b, mdParser, nil)

			content = string(html)

			// save to cache
			mdCache.Set(p.Content, content, cache.DefaultExpiration)
			log.Println("saved markdown to cache")
		}
	}

	return
}
