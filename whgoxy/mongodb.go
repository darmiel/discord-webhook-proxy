package whgoxy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strings"
	"time"
)

type SavedWebhook struct {
	UUID       string `json:"uuid"`
	WebhookUrl string `json:"webhook_url"`
	Data       bson.M `json:"data"`
}

var errorEmptyWebhookUrl error = errors.New("empty webhook url")

func (w *SavedWebhook) CheckValidity() (err error) {
	if len(w.WebhookUrl) <= 0 {
		return errorEmptyWebhookUrl
	}

	// TODO: Add more validity checks (discord format etc.)

	return nil
}

func (w *SavedWebhook) Send(param ...map[string]string) (err error) {
	// check if webhook is valid
	if err := w.CheckValidity(); err != nil {
		return err
	}

	// marshall data
	jsdb, err := json.Marshal(w.Data)
	if err != nil {
		return err
	}

	jsd := string(jsdb)

	// replace params in data
	if param != nil && len(param) >= 1 && len(param[0]) > 0 {
		for key, value := range param[0] {
			re := strings.NewReplacer(
				fmt.Sprintf("{{%s}}", key), value,
				fmt.Sprintf("{{ %s }}", key), value,

				fmt.Sprintf("{{ %s}}", key), value, // also \       / "faulty" \            /
				fmt.Sprintf("{{%s }}", key), value, //       replace            placeholders
			)

			jsd = re.Replace(jsd)
		}
	}

	// TODO: Debug. Remove me!
	log.Println("Sending json to discord:", jsd)

	// Send to discord
	reader := bytes.NewReader([]byte(jsd))
	req, err := http.NewRequest("POST", w.WebhookUrl, reader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return nil
}

func (w *SavedWebhook) Save(ctx context.Context, client *mongo.Client) (err error) {
	filter := bson.M{
		"uuid": w.UUID,
	}
	opts := options.Update().SetUpsert(true)
	update := bson.M{
		"$set": w,
	}
	_, err = client.Database(Opt.MongoDatabase).Collection("whgoxy").UpdateOne(ctx, filter, update, opts)
	return err
}

func buildApplyURI(opt *Options) (res string) {
	res = "mongodb+srv://"

	// add user
	res += opt.MongoAuthUser
	res += ":"

	// add pass
	res += opt.MongoAuthPass
	res += "@"

	// add host
	res += opt.MongoHost
	res += "/"

	// add db
	res += opt.MongoDatabase

	// add other params
	res += "?retryWrites=true&w=majority"

	return res
}

func InitMongoDatabase(opt *Options) {
	applyURI := buildApplyURI(opt)
	log.Println("ApplyURI:", applyURI)

	client, err := mongo.NewClient(options.Client().ApplyURI(applyURI))
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalln("Error while disconnecting:", err.Error())
		}
	}()
}
