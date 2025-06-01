package goredis

import (
	"github.com/go-redis/redis"
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"go_game_server/server/util"
)

// https://github.com/go-redis/redis

var ClientRedis *redis.Client

func InitRedis() {
	address := global.MyConfig.Read("redis", "address")
	password := global.MyConfig.Read("redis", "password")
	db := global.MyConfig.Read("redis", "db")
	poolsize := global.MyConfig.Read("redis", "poolsize")
	ClientRedis = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       int(util.ToInt(db)),
		PoolSize: int(util.ToInt(poolsize)),
	})

	//pong, err := ClientRedis.Ping().Result()
	//fmt.Println(pong, err) // Output: PONG <nil>
	logger.Log.Infof(">>>>>>>>>>>> redis启动成功,端口:%v \n\n", address)
}
