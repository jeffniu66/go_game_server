package global

import "sync"

// 命令行参数
var ServerPort *string  // 端口号
var Battle2Game *string // kafka队列名
var Game2Battle *string
var ServerNum int32

var PlayerGMList map[string]bool // 玩家开启的GM命令
var RobotBirthPointMap sync.Map
