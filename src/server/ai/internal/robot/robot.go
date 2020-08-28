package robot

import (
	"reflect"
	"server/conf"
	"server/models"
	"sync"
	"time"
)

var (
	uniqueId = 100000
 	callFun sync.Map
	robotList sync.Map
)

type Robot struct {
	CallCh      chan interface{}
	CloseCh     chan bool	//机器人关闭
	WorkEndCh 	chan bool //工作结束
	MatchTicker *time.Ticker
	GameInfo    map[string]interface{}
	*models.User
}

type Call struct {
	GameId int
	Msg interface{}
}

func (r *Robot) GetUniqueId() int  {
	uniqueId++
	return uniqueId
}

func (r *Robot) Create() *Robot {
	robot := new(Robot)
	callCh := make(chan interface{}, 10)
	closeCh:=make(chan bool, 1)
	workEndCh:=make(chan bool, 1)

	agent := new(robotAgent)
	agent.robotCallCh = callCh

	robot.User = &models.User{
		Uid : r.GetUniqueId(),
		Name: "robot",
		IsRobot:true,
	}
	robot.MatchTicker = time.NewTicker(conf.AiMatchTime)
	robot.CallCh = callCh
	robot.CloseCh = closeCh
	robot.WorkEndCh = workEndCh
	robot.User = robot.SetLoginInfo(robot.User, agent)
	robotList.LoadOrStore(robot.Uid, robot.CloseCh)
	return robot
	//r.Start()
}

func (r *Robot) Close()  {
	robotList.Delete(r.Uid)
	r.DeleteCache(r.Uid)
}

//接收消息返回调用函数
func (r *Robot) Call(message interface{}) {
	callFunMap, ok := callFun.Load(r.GameId)
	if !ok{
		return
	}

	callFunMap2 := callFunMap.(sync.Map)
	fun, ok := callFunMap2.Load(reflect.TypeOf(message))
	if !ok{
		return
	}

	fun.(func(interface{}, *Robot))(message, r)
}


//注册回调的消息调用的函数
func (r *Robot) RegisterCallMsg(gameId int, message interface{}, fun func(interface{}, *Robot))  {
	callFunMap, ok := callFun.Load(gameId)
	if !ok{
		callFunMap = sync.Map{}
	}
	callFunMap2 := callFunMap.(sync.Map)


	callFunMap2.Store(reflect.TypeOf(message), fun)
	callFun.Store(gameId, callFunMap2)
}

func (r *Robot) CloseAll()  {
	robotList.Range(func(key, value interface{}) bool {
		if ch,ok := value.(chan bool);ok{
			ch<-true
		}
		return true
	})
}