package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-web-app/dao/mysql"
	"go-web-app/logic"
	"go-web-app/models"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func CrontabSystem(c *gin.Context) {
	p := new(models.ParameCrontab)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误,直接返回响应

		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResopnseError(c, CodeServerApiType)
			return
		}
		ResponseErrorwithMsg(c, CodeServerApiType, removeTopStruct(errs.Translate(trans)))
		return
	}
	s, err := logic.Crond(p)
	if err != nil {
		zap.L().Error("ceontab with invalid param", zap.String("ParameOption", p.ParameOption), zap.Error(err))
		ResopnseError(c, CodeInvalidParam)

		return
	}
	//3.返回响应

	ResopnseSystemDataSuccess(c, s)
}
func DownloadHandler(c *gin.Context) {
	p := new(models.ParameCrontab)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误,直接返回响应

		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResopnseError(c, CodeServerApiType)
			return
		}
		ResponseErrorwithMsg(c, CodeServerApiType, removeTopStruct(errs.Translate(trans)))
		return
	}
	attchmentName := mysql.FileName(p.FileId)
	attchmentDir := mysql.FileDir(p.FileId)
	_, err := os.Open(attchmentDir)
	if err != nil {
		fmt.Println("文件获取失败", attchmentDir)
		fmt.Println(err)
		c.Redirect(http.StatusFound, "/404")
		return
	}

	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", attchmentName))
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.File(attchmentDir)
	return
}
