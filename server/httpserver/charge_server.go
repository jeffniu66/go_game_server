package httpserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_game_server/server/global"
	"go_game_server/server/util"
)

// 启动gin服务
func InitChargeGin() {
	start := global.MyConfig.Read("charge", "start")
	if start == "0" {
		return
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	ginChargeRouter(r)
	address := global.MyConfig.Read("charge", "address")
	r.Run(address)
}

// 开始添加所有http方法
func ginChargeRouter(router *gin.Engine) {
	fmt.Println("------------------ charge start success ---------------------")
	router.POST("/charge", postCharge)
}

func postCharge(c *gin.Context) {
	sign := util.GenMd5Sign(c.PostForm("uid"),
		c.PostForm("productID"),
		c.PostForm("usdAmount"),
		c.PostForm("actRate"),
		c.PostForm("orderNo"),
		global.MyConfig.Read("server", "sign_key"))

	fmt.Println("sign: ", sign)
}
