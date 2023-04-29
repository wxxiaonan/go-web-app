package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-web-app/dao/mysql"
	"go-web-app/logic"
	"go-web-app/models"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	//1.获取参数和校验参数
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误,直接返回响应
		zap.L().Error("SingUp with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResopnseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorwithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResopnseError(c, CodeUserExist)
			return
		}
		ResopnseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResopnseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	//获取请求参数及参数校验
	//业务逻辑处理
	//返回响应
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误,直接返回响应

		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResopnseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorwithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("SingUp with invalid param", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNoExist) {
			ResopnseError(c, CodeUserNotExist)
			return
		}
		ResopnseError(c, CodeInvalidPassword)

		return
	}
	//3.返回响应

	ResopnseSuccess(c, token)
}

func Systemdata(c *gin.Context) {
	p := new(models.ParamSystemGet)
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误,直接返回响应

		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResopnseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorwithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	s, err := logic.NetworkSentSpeed(p)
	if err != nil {
		zap.L().Error("SingUp with invalid param", zap.String("username", p.Type), zap.Error(err))
		ResopnseError(c, CodeServerApiType)

		return
	}
	//3.返回响应

	ResopnseSuccess(c, s)
}
