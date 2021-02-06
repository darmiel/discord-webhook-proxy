package discord

import (
	"context"
	"github.com/darmiel/whgoxy/internal/whgoxy/db/dbredis"
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
