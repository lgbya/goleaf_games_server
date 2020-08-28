package tictactoe

import (
	"github.com/chenhg5/collection"
	"github.com/name5566/leaf/gate"
	"github.com/shopspring/decimal"
	"math/rand"
	"reflect"
	"server/ai/internal/robot"
	"server/ai/internal/work"
	"server/define"
	"server/game"
	"server/msg"
	"time"
)

type Job struct {
	work.Work
}

func init() {
	job := new(Job)
	job.RegisterCallMsg(define.Tictactoe, &msg.S2C_StartGame{}, job.StartGame)
	job.RegisterCallMsg(define.Tictactoe, &msg.S2C_TictactoePlay{}, job.Play)
	job.RegisterCallMsg(define.Tictactoe, &msg.S2C_EndGame{}, job.EndGame)
}

func (j *Job)StartGame(message interface{}, r *robot.Robot)  {
	protocol := message.(*msg.S2C_StartGame)
	surplusList := []int{1,2,3,4,5,6,7,8,9}

	start := protocol.Start
	currentUid := start["currentUid"]
	surplusList = j.send(r.Uid, currentUid.(int), surplusList, *r.Agent )
	gameInfo := r.GameInfo
	gameInfo = make(map[string]interface{})
	gameInfo["surplus"] = surplusList
	r.GameInfo = gameInfo
}

func (j *Job)Play(message interface{}, r *robot.Robot)  {
	protocol := message.(*msg.S2C_TictactoePlay)
	gameInfo := r.GameInfo
	surplusList := gameInfo["surplus"].([]int)
	//先删除对手选择的
	surplusList = collection.Collect(surplusList).Reject(func(item, value interface{}) bool {
		return value.(decimal.Decimal).IntPart() == int64(protocol.Number)
	}).ToIntArray()
	//判断是否到机器人操作
	surplusList = j.send(r.Uid, protocol.CurrentUid, surplusList, *r.Agent )
	gameInfo["surplus"] = surplusList
	r.GameInfo = gameInfo
}

func (j *Job) send (robotUid int, currentUid int, surplusList []int, agent gate.Agent ) []int {
	if robotUid == currentUid {
		number := 0
		if len(surplusList) == 1 {
			number = surplusList[0]
			surplusList = []int{}
		}else{
			key := rand.Intn(len(surplusList)-1)
			number = surplusList[key]
			surplusList = append(surplusList[:key],surplusList[key+1:]...)
		}
		send := &msg.C2S_TictactoePlay{Number:number}
		go func() {
			<-time.After(3 * time.Second)
			game.ChanRPC.Go(reflect.TypeOf(send), send, agent)
		}()
	}
	return surplusList
}

func (j *Job)EndGame(protocol interface{}, r *robot.Robot)  {
	r.WorkEndCh<-true
}