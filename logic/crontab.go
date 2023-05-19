package logic

import (
	"go-web-app/dao/mysql"
	"go-web-app/models"
	"runtime"
)

func CrontabInit() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func Crond(p *models.ParameCrontab) (Reply interface{}, err error) {
	switch {
	case p.ParameOption == "add":
		Reply, err := mysql.CrontabAdd(p)
		if err != nil {
			return Reply, err
		}
		return Reply, err
	}
	return
}
