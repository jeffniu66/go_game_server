package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"go_game_server/server/util"
)

var DB *sql.DB

func InitMysql() {
	var err error
	address := global.MyConfig.Read("mysql", "address")
	username := global.MyConfig.Read("mysql", "username")
	password := global.MyConfig.Read("mysql", "password")
	database := global.MyConfig.Read("mysql", "database")
	address = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", username, password, address, database)
	DB, err = sql.Open("mysql", address)
	if err != nil {
		logger.Log.Errorln("mysql open fail!", err)
		return
	}

	err = DB.Ping()
	util.CheckErr(err)
	logger.Log.Infof(">>>>>>>>>>>> mysql启动成功,端口:%v \n\n", global.MyConfig.Read("mysql", "address"))
	DB.SetMaxOpenConns(120)
	DB.SetMaxIdleConns(60)
	InitDBPid()
}
