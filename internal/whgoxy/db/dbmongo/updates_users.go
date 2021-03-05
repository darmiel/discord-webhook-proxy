package dbmongo

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func (mdb *mongoDatabase) SaveDiscordUser(dcu *discord.DiscordUser) (err error) {
	log.Println("💾 Saving user", dcu.GetFullName(), "...")

	filter := bson.M{"id": dcu.UserID}
	update := bson.M{"$set": dcu}

	updateOpts := options.Update().SetUpsert(true)

	_, err = mdb.userCollection().UpdateOne(mdb.context, filter, update, updateOpts)
	return
}
