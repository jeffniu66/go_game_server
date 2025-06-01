package game

import (
	"go_game_server/proto3"
	"go_game_server/server/db"
	"go_game_server/server/global"
)

// 发送全服邮件
func SendFullServerMail(title, content, reward string) {
	userIds := db.GetUserIds()
	if len(userIds) == 0 {
		return
	}
	for _, v := range userIds {
		userMail := &db.UserMail{
			UserID:     v,
			Title:      title,
			Content:    content,
			AnnexItems: reward,
			DetailType: int32(proto3.DetailTypeEnum_full_server_text),
		}
		db.SaveMail(userMail)
		player := global.GloInstance.GetPlayer(v)
		if player == nil {
			continue
		}
		p := player.(*Player)
		p.SendRedPoint(proto3.RedPointEnum_mail_red, []int32{1})
	}
}
