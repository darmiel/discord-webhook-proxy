package dbmongo

import (
	"errors"
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/cms"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

// FindCMSPage ...
func (mdb *mongoDatabase) FindCMSPage(url string) (page *cms.CMSPage, err error) {
	// check cache
	if p, found := db.CMSCache.Get("page::" + url); found {
		log.Println("CMS page found in cache:", p)
		a, b := p.(*cms.CMSPage)
		if b {
			log.Println("succ:", a)
		} else {
			log.Println("err:", a)
		}

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

func (mdb *mongoDatabase) FindAllCMSPages() (pages []*cms.CMSPage, err error) {

	// check cache
	if p, found := db.CMSCache.Get("*::all"); found {
		return p.([]*cms.CMSPage), nil
	}

	filter := bson.M{}

	res, err := mdb.cmsCollection().Find(mdb.context, filter)
	if err != nil {
		return nil, err
	}

	for res.Next(mdb.context) {
		var page *cms.CMSPage
		if err := res.Decode(&page); err != nil {
			return nil, err
		}

		//goland:noinspection GoNilness
		if page == nil {
			return nil, errors.New("a cms page was nil")
		}

		pages = append(pages, page)

		// update cache for specific page
		db.CMSCache.Set("page::"+page.URL, page, cache.DefaultExpiration)
	}

	// update cache
	db.CMSCache.Set("*::all", pages, cache.DefaultExpiration)

	return pages, nil
}

func (mdb *mongoDatabase) FindAllLinks() (links []*cms.CMSLink, err error) {
	// check cache
	if p, found := db.CMSCache.Get("link::*::all"); found {
		return p.([]*cms.CMSLink), nil
	}

	filter := bson.M{}
	res, err := mdb.linkCollection().Find(mdb.context, filter)
	if err != nil {
		return nil, err
	}

	for res.Next(mdb.context) {
		var link *cms.CMSLink
		if err := res.Decode(&link); err != nil {
			return nil, err
		}

		//goland:noinspection GoNilness
		if link == nil {
			return nil, errors.New("a cms link was nil")
		}

		links = append(links, link)
	}

	// update cache
	db.CMSCache.Set("links::*::all", links, cache.DefaultExpiration)
	return links, nil
}
