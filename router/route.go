package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web-app/controller"
	"go-web-app/dao/mysql"
	"go-web-app/logger"
	"go-web-app/middlewares"
	"go-web-app/models"
	"net/http"
)

func Setup(mode, ClientUrl string, size int64, savedir string) *gin.Engine {
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middlewares.Cors(ClientUrl))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.MaxMultipartMemory = size << 20
	//注册业务路由1·1··
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
	r.POST("/crontab", controller.CrontabSystem)
	r.POST("/download", controller.DownloadHandler)
	r.POST("/upload", func(ctx *gin.Context) {
		forms, err := ctx.MultipartForm()
		if err != nil {
			fmt.Println("error", err)
		}
		files := forms.File["file"]
		for _, v := range files {
			filelog := &models.Filelog{
				FileName: v.Filename,
				FileSize: v.Size,
				FileDir:  savedir + v.Filename,
			}
			fmt.Println(filelog)
			if err := ctx.SaveUploadedFile(v, fmt.Sprintf("%s%s", savedir, v.Filename)); err != nil {
				ctx.String(http.StatusBadRequest, fmt.Sprintf("upload err %s", err.Error()))
			}
			err = mysql.FileLogAdd(filelog)
		}
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
