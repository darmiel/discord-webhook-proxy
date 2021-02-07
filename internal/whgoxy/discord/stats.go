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
	CallsGlobal       uint64
	Calls60           uint64
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

	stats.CallsGlobal, _ = dbredis.Get(w.UserID, w.UID, dbredis.KeyCallGlobalCount).Uint64()
	stats.Calls60, _ = dbredis.Get(w.UserID, w.UID, dbredis.KeyCallPerMinuteCount).Uint64()

	return
}

func (w *Webhook) AddCallStats() {
	redis := dbredis.GlobalRedis

	/// global call count
	redis.Incr(
		context.TODO(),
		dbredis.GetKey(w.UserID, w.UID, dbredis.KeyCallGlobalCount),
	)
	///

	/// calls per minute
	callMinuteKey := dbredis.GetKey(w.UserID, w.UID, dbredis.KeyCallPerMinuteCount)
	expire, _ := redis.Exists(
		context.TODO(),
		callMinuteKey,
	).Result()

	// increment
	redis.Incr(
		context.TODO(),
		callMinuteKey,
	)

	// expire?
	if expire == 0 {
		redis.Expire(
			context.TODO(),
			callMinuteKey,
			60*time.Second,
		)
	}
}
