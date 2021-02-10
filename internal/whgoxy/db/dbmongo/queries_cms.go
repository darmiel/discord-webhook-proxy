package dbmongo

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/cms"
	"go.mongodb.org/mongo-driver/bson"
)

// FindCMSPage ...
func (mdb *mongoDatabase) FindCMSPage(url string) (page *cms.CMSPage, err error) {
	// check cache
	if p, found := db.CMSCache.Get(url); found {
		return p.(*cms.CMSPage), nil
	}

	filter := &bson.M{
		"url": url,
	}

	res := mdb.cmsCollection().FindOne(mdb.context, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	page = &cms.CMSPage{}
	err = res.Decode(page)

	return
}