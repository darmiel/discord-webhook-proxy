package dbmongo

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/cms"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (mdb *mongoDatabase) SaveCMSPage(page cms.CMSPage) (err error) {
	filter := page.CreateFilter()
	update := bson.M{"$set": page}

	updateOpts := options.Update().SetUpsert(true)

	_, err = mdb.cmsCollection().UpdateOne(mdb.context, filter, update, updateOpts)

	// save to cache
	if err == nil {
		db.CMSCache.Set(page.URL, &page, cache.DefaultExpiration)
	}

	return
}

func (mdb *mongoDatabase) DeleteCMSPage(page cms.CMSPage) (err error) {
	filter := page.CreateFilter()

	_, err = mdb.cmsCollection().DeleteOne(mdb.context, filter)

	// remove from cache
	if err == nil {
		db.CMSCache.Delete(page.URL)
	}

	return
}
