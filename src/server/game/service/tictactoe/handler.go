package tictactoe

import (
	"reflect"
	"server/game/internal"
	"server/models"
	"server/msg"
)

var lMatch = make(map[int]int)

type Tictactoe struct {
	Info map[int]int
}

func init() {
}

func handler(m interface{}, h interface{})  {
	internal.GetSkeleton().RegisterChanRPC(reflect.TypeOf(m), h)
}

func (t *Tictactoe) MatchPlayer(user *models.User, protocol interface{})  {

}

func (t *Tictactoe) CancelMatch(user *models.User, protocol interface{})  {
	delete(lMatch, user.Uid)
	(*user.Agent).WriteMsg(&msg.S2C_MatchPlayer{})
}

func (t *Tictactoe) StartGame(room *models.Room)  {

}

func (t *Tictactoe) handlerMoraPlaying(args []interface{}) {

}

func (t *Tictactoe) endGame(room *models.Room)  {

}