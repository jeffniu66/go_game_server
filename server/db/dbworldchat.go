package db

import (
	"go_game_server/server/include"
	"go_game_server/server/util"
)

type WorldChat = include.WorldChat

func SaveWorldChat(userID, userPhoto int32, username, chatData string) {
	intSQL := "INSERT INTO t_world_chat(user_id, username, user_photo, chat_data, create_time)VALUES(?, ?, ?, ?, ?)"
	ExecDB(UserDBType, userID, intSQL, userID, username, userPhoto, chatData, util.UnixTime())
	return
}

func ListWorldChat() []WorldChat {
	return nil
}
