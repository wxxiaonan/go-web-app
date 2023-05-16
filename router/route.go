package router

import (
	"github.com/gin-gonic/gin"
	"go-web-app/controller"
	"go-web-app/logger"
	"go-web-app/middlewares"
	"net/http"
)

func Setup(mode, ClientUrl string) *gin.Engine {
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middlewares.Cors(ClientUrl))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//注册业务路由
	r.POST("/signup", controller.SignUpHandler)
	r.POST("/login", controller.LoginHandler)

	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		c.String(http.StatusOK, "pong")

	})
	r.POST("/systemview", controller.Systemdata)
	r.POST("/hostlistdata", controller.Hostlistata)
	r.POST("/statisticsdata", controller.Statisticsdata)
	r.POST("/alarmsetting", controller.Alarmsetting)
	r.POST("/clientdata", controller.Clientdata)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
