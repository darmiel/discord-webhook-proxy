package cms

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type CMSPageMeta struct {
	Title         string    `bson:"title" json:"title"`
	CreatorUserID string    `bson:"creator_user_id" json:"creator_user_id"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
}

type CMSPagePreferences struct {
	AuthorVisible    bool `bson:"author_visible" json:"author_visible"`
	Dynamic          bool `bson:"dynamic" json:"dynamic"` // Does the page include GoLang-Template content?
	URLCaseSensitive bool `bson:"url_case_sensitive" json:"url_case_sensitive"`
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
