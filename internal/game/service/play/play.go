package play

import (
	"server/internal/model"
)

type Play interface {
	Start(*model.Room, ...interface{}) (map[string]interface{}, *model.Room)  //每个游戏如果需要单独处理开始的钩子
	Continue(*model.User, *model.Room, ...interface{}) map[string]interface{} //每个游戏如果需要单独处理开始的钩子
	Run(*model.Call)
}