package logic

import (
	"fmt"
	"go-web-app/dao/mysql"
	"go-web-app/models"
	"go-web-app/pkg/jwt"
	"go-web-app/pkg/snowflake"
)

// 存放业务逻辑代码

func SignUp(p *models.ParamSignUp) (err error) {

	//1.判断用户可用性
	err = mysql.CheckUserByUsername(p.Username)
	if err != nil {
		//数据库查询出错
		return err
	}

	//2.生成UID
	userId := snowflake.GenID()
	//构造一个user示例
	user := &models.User{
		UserId:   userId,
		Username: p.Username,
		Password: p.Password,
	}
	//3.用户数据入库
	return mysql.InsertUser(user)

}

func Login(p *models.ParamLogin) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//传递的是指针，可以获取到user.UserId
	if err := mysql.Login(user); err != nil {
		return "", err

	}
	fmt.Println(err)

	//生成JWTtoken
	return jwt.GenToken(user.UserId, user.Username)

}
