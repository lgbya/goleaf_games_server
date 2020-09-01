package more

import (
	"server/internal/ai/robot"
	"server/internal/ai/work"
	"server/internal/common/define"
	"server/internal/game"
	"server/internal/protocol"

	"math/rand"
	"reflect"
)

type Job struct {
	*work.Work
}

func init() {
	job := &Job{}
	robot.RegisterCallMsg(define.More, &protocol.S2C_StartGame{}, job.StartGame)
	robot.RegisterCallMsg(define.More, &protocol.S2C_EndGame{}, job.EndGame)
}

func (j Job) StartGame(message interface{}, r *robot.Robot) {
	send := &protocol.C2S_MoraPlay{Ply: rand.Intn(2) + 1}
	game.ChanRPC.Go(reflect.TypeOf(send), send, *r.User.Agent)
}

func (j Job) EndGame(message interface{}, r *robot.Robot) {
	r.WorkEndCh <- true
}
