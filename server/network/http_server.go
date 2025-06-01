package network

import (
	"go_game_server/server/db"
	"go_game_server/server/game"
	"go_game_server/server/global"
	"go_game_server/server/handler"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
	"net/http"
	"net/http/pprof"
	"net/url"

	"github.com/gin-gonic/gin"
)

func StartHttpServer() {
	addr := global.MyConfig.Read("server", "httpaddr")
	logger.Log.Infof("http server started on %s\n", addr)
	router := gin.Default()
	router.GET("/gm", handleGm)
	router.GET("/serverconfig/update", handleServerConfig)
	router.GET("/tableconfig/update", handleTableConfig)
	router.GET("/msgsize", handleMsgSize)
	router.GET("/changeServerStatus", changeServerStatus)
	router.POST("/sendFullServerMail", sendFullServerMail)

	routeRegisterPprof(router)
	router.Run(addr)
}

func routeRegisterPprof(rg *gin.Engine) {
	prefix := "/debug/pprof"

	prefixRouter := rg.Group(prefix)
	{
		prefixRouter.GET("/", pprofHandler(pprof.Index))
		prefixRouter.GET("/cmdline", pprofHandler(pprof.Cmdline))
		prefixRouter.GET("/profile", pprofHandler(pprof.Profile))
		prefixRouter.POST("/symbol", pprofHandler(pprof.Symbol))
		prefixRouter.GET("/symbol", pprofHandler(pprof.Symbol))
		prefixRouter.GET("/trace", pprofHandler(pprof.Trace))
		prefixRouter.GET("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
		prefixRouter.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
		prefixRouter.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
		prefixRouter.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
		prefixRouter.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
		prefixRouter.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
	}
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	handler := http.HandlerFunc(h)
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func sendFullServerMail(c *gin.Context) {
	json := make(map[string]string)
	c.BindJSON(&json)
	logger.Log.Infof("%v", &json)
	title := json["title"]
	content := json["content"]
	reward := json["reward"]
	go game.SendFullServerMail(title, content, reward)
	c.String(200, "ok")
}

func handleGm(c *gin.Context) {
	gm := c.Query("gm")
	msg := c.Query("msg")
	gmUser := &game.Player{Attr: &db.Attr{}}
	gmUser.Attr.Username = gm
	data := ""
	msg, _ = url.QueryUnescape(msg)
	handler.GmCommand(gmUser, msg, &data)
	c.String(200, data)
}

func handleServerConfig(c *gin.Context) {
	logger.Log.Info("更新ServerConfig")
	global.MyConfig.InitConfig("./config/pro_server.config")
	c.String(200, "ok")
}

func handleTableConfig(c *gin.Context) {
	logger.Log.Info("更新tableConfigs")
	// 读策划表数据
	if err := tableconfig.ReadTable(); err != nil {
		logger.Log.Errorf("read table failed, err:%v", err)
		return
	}
	// 刷新定时任务
	game.FreshTickerManageData()
	c.String(200, "ok")
}

func handleMsgSize(c *gin.Context) {
	logger.Log.Info("get msg size")
	// 读策划表数据
	data := sendMsgSize
	c.String(200, util.ToInt64Str(data)+"B, "+util.ToInt64Str(data/1024)+"KB")
}

func changeServerStatus(c *gin.Context) {
	if ServerStatus {
		ServerStatus = false
	} else {
		ServerStatus = true
	}
	c.String(200, "ok")
}
