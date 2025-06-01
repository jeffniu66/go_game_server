package util

import (
	"database/sql"
	"encoding/json"
	"github.com/go-redis/redis"
	"runtime/debug"
	"go_game_server/server/logger"
)

func CheckErr(err error) {
	if err == nil {
		return
	}
	switch {
	case err == sql.ErrNoRows:
	case err == redis.Nil:
		logger.Log.Errorln("redis 有不存在的Key--------- err = ", err, string(debug.Stack()))
	default:
		logger.Log.Errorln("------------------- err = ", err, string(debug.Stack()))
	}
}

func CheckRedisErr(key string, err error) {
	if err == nil {
		return
	}
	switch {
	case err == sql.ErrNoRows:
	case err == redis.Nil:
		logger.Log.Errorln("redis 有不存在的Key--------- err = ", key, err, string(debug.Stack()))
	default:
		logger.Log.Errorln("------------------- err = ", key, err, string(debug.Stack()))
	}
}

func CheckRedisErrs(keys []string, err error) {
	if err == nil {
		return
	}
	switch {
	case err == sql.ErrNoRows:
	case err == redis.Nil:
		logger.Log.Errorln("redis 有不存在的Key--------- err = ", keys, err, string(debug.Stack()))
	default:
		logger.Log.Errorln("------------------- err = ", keys, err, string(debug.Stack()))
	}
}

func CheckJsonMarshalErr(v interface{}) (str []byte) {
	str, err := json.Marshal(v)
	CheckErr(err)
	return
}
