package play

import (
	"server/internal/common/define"
	"server/internal/game/service/play/mora"
	"server/internal/game/service/play/tictactoe"
	"server/internal/model"
)

type Play interface {
	//MatchPlayer(*models.User, ...interface{}) 	//匹配玩家
	Start(*model.Room, ...interface{}) (map[string]interface{}, *model.Room)  //每个游戏如果需要单独处理开始的钩子
	Continue(*model.User, *model.Room, ...interface{}) map[string]interface{} //每个游戏如果需要单独处理开始的钩子
	Run(*model.Call)
}

func New(gameId int) (Play, bool) {
	serviceList := map[int]Play{
		define.More:      new(mora.Mode),
		define.Tictactoe: new(tictactoe.Mode),
	}
	service, ok := serviceList[gameId]
	return service, ok
}

func AllGameId() []int {
	return []int{
		define.More,
		define.Tictactoe,
	}
}
