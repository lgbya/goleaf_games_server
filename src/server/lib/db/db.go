package db

import (
	"github.com/jinzhu/gorm"
	"github.com/name5566/leaf/log"
	_ "github.com/go-sql-driver/mysql"
	"server/conf"
)

var mysqlDb *gorm.DB

func init() {
	cfg := conf.Server.SqlSrv
	args := cfg.Username + ":" + cfg.Password + "@tcp("+ cfg.Host + ")/"+ cfg.Database +"?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql",  args)
	if err != nil {
		log.Debug("连接数据库失败%v, %v", args, err)
		panic(err)
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(15)
	db.SingularTable(true)
	//log.Release("===========连接数据库成功=========== \n")
	mysqlDb = db
}

func New() *gorm.DB {
	return mysqlDb
}