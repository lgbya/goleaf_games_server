package play

import (
	"server/define"
	"server/game/service/play/mora"
	"server/game/service/play/tictactoe"
	"server/models"
)


type Play interface {
	//MatchPlayer(*models.User, ...interface{}) 	//匹配玩家
	Start(*models.Room, ...interface{})	 (map[string]interface{}, *models.Room)//每个游戏如果需要单独处理开始的钩子
	Continue(*models.User, *models.Room, ...interface{}) map[string]interface{}	//每个游戏如果需要单独处理开始的钩子
	Run(*models.Call)
}


func New(gameId int ) (Play, bool)  {
	serviceList := map[int]Play{
		define.More:      new(mora.Mode),
		define.Tictactoe: new(tictactoe.Mode),
	}
	service, ok := serviceList[gameId]
	return service, ok
}

func AllGameId()[]int{
	return []int{
		define.More,
		define.Tictactoe,
	}
}