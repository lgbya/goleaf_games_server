package service

import (
	"server/internal/common/define"
	"server/internal/game/service/play"
	"server/internal/game/service/play/mora"
	"server/internal/game/service/play/tictactoe"
)

func NewPlay(gameId int) (play.Play, bool) {
	playServiceList := map[int]play.Play{
		define.More:      new(mora.Mode),
		define.Tictactoe: new(tictactoe.Mode),
	}
	playService, ok := playServiceList[gameId]
	return playService, ok
}

func AllGameId() []int {
	return []int{
		define.More,
		define.Tictactoe,
	}
}
