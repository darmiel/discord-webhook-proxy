package dbmongo

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/cms"
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
		db.CMSCache.Delete("*::all")
		db.CMSCache.Delete("page::" + page.URL)
	}

	return
}

func (mdb *mongoDatabase) DeleteCMSPage(page cms.CMSPage) (err error) {
	filter := page.CreateFilter()

	_, err = mdb.cmsCollection().DeleteOne(mdb.context, filter)

	// remove from cache
	if err == nil {
		db.CMSCache.Delete("*::all")
		db.CMSCache.Delete("page::" + page.URL)
	}

	return
}

func (mdb *mongoDatabase) SaveCMSLink(link *cms.CMSLink) (err error) {
	filter := bson.M{
		"name": link.Name,
	}
	update := bson.M{
		"$set": link,
	}
	_, err = mdb.linkCollection().UpdateOne(mdb.context, filter, update, options.Update().SetUpsert(true))
	// remove from cache
	if err == nil {
		db.CMSCache.Delete("link::*::all")
	}
	return
}

func (mdb *mongoDatabase) DeleteCMSLink(link *cms.CMSLink) (err error) {
	filter := bson.M{
		"name": link.Name,
	}
	_, err = mdb.linkCollection().DeleteOne(mdb.context, filter)
	// remove from cache
	if err == nil {
		db.CMSCache.Delete("link::*::all")
	}
	return
}
