package service

import (
	"server/game/service/mora"
	"server/game/service/tictactoe"
	"server/models"
)

type GameService interface {
	//MatchPlayer(*models.User, ...interface{}) 	//匹配玩家
	Start(*models.Room, ...interface{})	 (map[string]interface{}, *models.Room)//每个游戏如果需要单独处理开始的钩子
	Continue(*models.User, *models.Room, ...interface{}) map[string]interface{}	//每个游戏如果需要单独处理开始的钩子
}


func NewGameService(gameId int ) (GameService, bool)  {
	serviceList := map[int]GameService{
		1001 : new(mora.Mode),
		1002 : new(tictactoe.Mode),
	}
	service, ok := serviceList[gameId]
	return service, ok
}