package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		viper.GetString("database.mysql.user"),
		viper.GetString("database.mysql.password"),
		viper.GetString("database.mysql.host"),
		viper.GetInt("database.mysql.port"),
		viper.GetString("database.mysql.dbname"))
	//初始化全局的DB对象
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect to db failed", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(viper.GetInt("databases.mysql.max_conns"))
	db.SetMaxIdleConns(viper.GetInt("databases.mysql.max_idle_conns"))
	return

}

func Close() {
	_ = db.Close()
}
