package cms

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type CMSPageMeta struct {
	Title         string
	CreatorUserID string
	CreatedAt     time.Time
}

type CMSPagePreferences struct {
	AuthorVisible    bool
	Dynamic          bool // Does the page include GoLang-Template content?
	URLCaseSensitive bool
}

type CMSPageUpdate struct {
	UpdatedAt     time.Time // Time of update
	UpdaterUserID string    // UserID of updater, -1 = System
}

type CMSPage struct {
	URL         string
	Meta        CMSPageMeta
	Updates     []CMSPageUpdate
	Preferences CMSPagePreferences
	Content     string
}

func (p *CMSPage) CreateFilter() *bson.M {
	return &bson.M{
		"url": p.URL,
	}
}
