package db

import (
	"server/internal/common/conf"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/name5566/leaf/log"
)

var _db *gorm.DB

func init() {
	config := conf.Get().DB
	args := config.Username + ":" + config.Password + "@tcp(" + config.Host + ")/" + config.Database + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", args)
	if err != nil {
		log.Debug("连接数据库失败%v, %v", args, err)
		panic(err)
	}
	db.DB().SetMaxIdleConns(config.MaxIdle)
	db.DB().SetMaxOpenConns(config.MaxOpen)
	db.SingularTable(true)
	//log.Release("===========连接数据库成功=========== \n")
	_db = db
}

func New() *gorm.DB {
	return _db
}
