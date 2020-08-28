package more

import (
	"math/rand"
	"reflect"
	"server/ai/internal/robot"
	"server/ai/internal/work"
	"server/define"
	"server/game"
	"server/msg"
)

type Job struct {
	work.Work
}

func init() {
	job := new(Job)
	job.RegisterCallMsg(define.More, &msg.S2C_StartGame{}, job.StartGame)
	job.RegisterCallMsg(define.More, &msg.S2C_EndGame{}, job.EndGame)
}


func (j *Job)StartGame(message interface{}, r *robot.Robot)  {
	send := &msg.C2S_MoraPlay{Ply: rand.Intn(2)+1}
	game.ChanRPC.Go(reflect.TypeOf(send), send, *r.Agent)
}

func (j *Job)EndGame(message interface{}, r *robot.Robot)  {
	r.WorkEndCh<-true
}