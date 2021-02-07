package discord

import (
	"context"
	"github.com/darmiel/whgoxy/internal/whgoxy/db/dbredis"
	"time"
)

func (w *Webhook) AddError(err error, json string) (reserr error) {
	// increment error count
	redis := dbredis.GlobalRedis
	reserr = redis.Incr(
		context.TODO(),
		dbredis.GetKey(w.UserID, w.UID, dbredis.KeyErrorCount),
	).Err()
	if reserr != nil {
		return
	}

	// set error message
	reserr = redis.Set(
		context.TODO(),
		dbredis.GetKey(w.UserID, w.UID, dbredis.KeyErrorMessage),
		err.Error(),
		7*24*time.Hour,
	).Err()
	if reserr != nil {
		return
	}

	// set json
	reserr = redis.Set(
		context.TODO(),
		dbredis.GetKey(w.UserID, w.UID, dbredis.KeyErrorJson),
		json,
		7*24*time.Hour,
	).Err()
	if reserr != nil {
		return
	}

	return nil
}

func (w *Webhook) AddSuccess() (reserr error) {
	// increment success count
	redis := dbredis.GlobalRedis
	return redis.Incr(
		context.TODO(),
		dbredis.GetKey(w.UserID, w.UID, dbredis.KeySuccessCount),
	).Err()
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

	// successful calls
	stats.SuccessfulCalls, _ = dbredis.Get(w.UserID, w.UID, dbredis.KeySuccessCount).Uint64()
	stats.UnsuccessfulCalls, _ = dbredis.Get(w.UserID, w.UID, dbredis.KeyErrorCount).Uint64()
	stats.LastErrorMessage, _ = dbredis.Get(w.UserID, w.UID, dbredis.KeyErrorMessage).Result()
	stats.LastErrorSentJson, _ = dbredis.Get(w.UserID, w.UID, dbredis.KeyErrorJson).Result()

	return
}
