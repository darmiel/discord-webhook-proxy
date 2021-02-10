package dbmongo

import (
	"errors"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/darmiel/whgoxy/internal/whgoxy/http/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func (mdb *mongoDatabase) FindDiscordUser(userID string) (dcu *discord.DiscordUser, err error) {
	// check cache
	if u, found := auth.AuthUserCache.Get(userID); found {
		return u.(*auth.User).DiscordUser, nil
	}

	filter := bson.M{
		"user_id": userID,
	}

	res := mdb.userCollection().FindOne(mdb.context, filter, options.FindOne())
	if res.Err() != nil {
		return nil, res.Err()
	}

	dcu = &discord.DiscordUser{}
	err = res.Decode(dcu)

	// repair user
	if err == nil && dcu.Repair() {
		log.Println("🔨 Repaired user", userID, "(", dcu.GetFullName(), ")")
		if e := mdb.SaveDiscordUser(dcu); e != nil {
			err = errors.New("error saving repaired user: " + e.Error())
		}
	}

	return
}