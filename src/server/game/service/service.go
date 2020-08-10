package service

import (
	"server/game/service/mora"
	"server/models"
)

type GameService interface {
	MatchPlayer(*models.User, interface{}) 	//匹配玩家
	CancelMatch(*models.User, interface{}) 	//匹配玩家
	StartGame(*models.Room)		//每个游戏如果需要单独处理开始的钩子
}



func  NewGameService(gameId int ) (GameService, bool)  {
	mapInt2GameService := map[int]GameService{
		1001 : new(mora.Mora),
	}
	service, ok := mapInt2GameService[gameId]
	return service, ok
}