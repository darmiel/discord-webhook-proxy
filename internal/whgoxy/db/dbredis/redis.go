package dbredis

import (
	"context"
	"github.com/darmiel/whgoxy/internal/whgoxy/config"
	"github.com/go-redis/redis/v8"
	"log"
)

func NewClient(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

var GlobalRedis *redis.Client

const (
	KeySuccessCount = iota
	KeyErrorCount
	KeyErrorMessage
	KeyErrorJson
	KeyCallGlobalCount
	KeyCallPerMinuteCount
)

func Get(userID, uid string, keyType int) *redis.StringCmd {
	key := GetKey(userID, uid, keyType)
	if key == "" {
		log.Println("[Redis] WARN: key", keyType, "resulted in an empty key")
		return &redis.StringCmd{}
	}
	return GlobalRedis.Get(context.TODO(), key)
}

func GetKey(userID, uid string, keyType int) (res string) {
	res = "whgoxy::stats::us:" + userID + "::ui:" + uid

	switch keyType {
	case KeySuccessCount:
		res += "::success:count"
		break
	case KeyErrorCount:
		res += "::error:count"
		break
	case KeyErrorMessage:
		res += "::error:msg"
		break
	case KeyErrorJson:
		res += "::error:json"
		break
	case KeyCallGlobalCount:
		res += "::calls:g"
		break
	case KeyCallPerMinuteCount:
		res += "::calls:60"
		break
	default:
		return ""
	}

	return
}
