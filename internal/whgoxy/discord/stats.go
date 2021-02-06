package discord

import (
	"context"
	"github.com/darmiel/whgoxy/internal/whgoxy/db/dbredis"
	"log"
	"time"
)

func (w *Webhook) AddError(err error, json string) (reserr error) {
	key := "whgoxy::stats::" + w.GetID() + "::error::"

	// increment error count
	redis := dbredis.GlobalRedis
	reserr = redis.Incr(
		context.TODO(),
		key+"count",
	).Err()
	if reserr != nil {
		return
	}

	// set error message
	reserr = redis.Set(
		context.TODO(),
		key+"msg",
		err.Error(),
		7*24*time.Hour,
	).Err()
	if reserr != nil {
		return
	}

	// set json
	reserr = redis.Set(
		context.TODO(),
		key+"json",
		json,
		7*24*time.Hour,
	).Err()
	if reserr != nil {
		return
	}

	return nil
}

func (w *Webhook) AddSuccess() (reserr error) {
	key := "whgoxy::stats::" + w.GetID() + "::success::count"

	// increment success count
	redis := dbredis.GlobalRedis
	reserr = redis.Incr(
		context.TODO(),
		key,
	).Err()

	return
}

type WebhookStats struct {
	SuccessfulCalls   uint64
	UnsuccessfulCalls uint64
	LastErrorMessage  string
	LastErrorSentJson string
}

func (w *Webhook) GetStats() (stats *WebhookStats) {
	stats = &WebhookStats{
		SuccessfulCalls:   0,
		UnsuccessfulCalls: 0,
		LastErrorMessage:  "",
		LastErrorSentJson: "",
	}

	redis := dbredis.GlobalRedis
	key := "whgoxy::stats::" + w.GetID() + "::"
	log.Println("key:", key)
	log.Println("redis:", redis)
	log.Println(redis.Get(context.TODO(), "heartbeat"))

	// successful calls
	stats.SuccessfulCalls, _ = redis.Get(context.TODO(), key+"success::count").Uint64()
	stats.UnsuccessfulCalls, _ = redis.Get(context.TODO(), key+"error::count").Uint64()

	stats.LastErrorMessage, _ = redis.Get(context.TODO(), key+"error::msg").Result()
	stats.LastErrorSentJson, _ = redis.Get(context.TODO(), key+"error::json").Result()

	return
}
