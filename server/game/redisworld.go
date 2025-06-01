package game

import (
	"encoding/json"
	"errors"
	"go_game_server/server/goredis"
	"go_game_server/server/logger"
	"go_game_server/server/util"
)

// 玩家活跃数据
type RedisUser struct {
	LordName string // 领主名字
	Domain   int32  // 领地个数
	Country  int32  // 所属国家
	Power    int32  // 势力值
}

// ----------------------------------------- 大地图上的活跃数据 --------------------------------------------------

func GetRedisUserKey(userID int32) string {
	return "user" + util.ToStr(userID)
}

func GetRedisUserField(userID int32, field string) (string, error) {
	return GetHashStrField(GetRedisUserKey(userID), field)
}

/********************************************************************************
**	redis 通用接口
****************************************************************************** */
// 是否存在key
func ExistKey(keys ...string) bool {
	val, err := goredis.ClientRedis.Exists(keys...).Result()
	if err != nil {
		logger.Log.Warnln("ExistKey err!", err)
		return false
	}

	return int64(len(keys)) == val
}

func SetHashStrField(key, field string, value interface{}) bool {
	_, err := goredis.ClientRedis.HSet(key, field, value).Result()
	util.CheckErr(err)
	return err == nil
}

func SetHashJsonValue(key, field string, value interface{}) bool {
	v, err := json.Marshal(value)
	util.CheckErr(err)
	if err != nil {
		return false
	}

	return SetHashStrField(key, field, v)
}

func SetHashFields(key string, fields map[string]interface{}) bool {
	_, err := goredis.ClientRedis.HMSet(key, fields).Result()
	util.CheckErr(err)
	return err == nil
}

func GetHashFields(key string, fields ...string) ([]interface{}, error) {
	if goredis.ClientRedis.Exists(key).Val() == 0 {
		return nil, errors.New("key not exists")
	}
	values, err := goredis.ClientRedis.HMGet(key, fields...).Result()
	util.CheckErr(err)
	return values, err
}

func GetHashStrField(key, field string) (string, error) {
	return goredis.ClientRedis.HGet(key, field).Result()
}

func GetHash(key string) (map[string]string, error) {
	values, err := goredis.ClientRedis.HGetAll(key).Result()
	if err != nil {
		logger.Log.Warnf("GetHash key:%v, err:%v", key, err)
		return nil, err
	}

	return values, err
}

func ListRPush(key string, val interface{}) {
	_, err := goredis.ClientRedis.RPush(key, val).Result()
	util.CheckErr(err)
}

func ListLPush(key string, val interface{}) {
	_, err := goredis.ClientRedis.LPush(key, val).Result()
	util.CheckErr(err)
}

func ListLPop(key string) {
	_, err := goredis.ClientRedis.LPop(key).Result()
	util.CheckErr(err)
}

func ListRPop(key string) {
	_, err := goredis.ClientRedis.RPop(key).Result()
	util.CheckErr(err)
}

func ListLen(key string) int64 {
	val, err := goredis.ClientRedis.LLen(key).Result()
	util.CheckErr(err)
	return val
}

func SetString(key string, val interface{}) {
	_, err := goredis.ClientRedis.Set(key, val, 0).Result()
	util.CheckErr(err)
}

func GetString(key string) string {
	val, err := goredis.ClientRedis.Get(key).Result()
	util.CheckErr(err)
	return val
}

func GetHashAllMapStr(key string) map[string]string {
	values, err := goredis.ClientRedis.HGetAll(key).Result()
	if err != nil {
		logger.Log.Warnf("GetHashAllMapStr key:%v, err:%v", key, err)
		return nil
	}

	return values
}

func GetRedisUserByKey(userKey string) *RedisUser {
	fields := GetHashAllMapStr(userKey)
	if fields == nil || len(fields) == 0 {
		return nil
	}
	redisUser := &RedisUser{}
	return redisUser
}
