package robot

import (
	"server/internal/common/conf"
	"server/internal/model"

	"math/rand"
	"reflect"
	"sync"
	"time"
)

var (
	_robotId   = 100000
	_callFun   sync.Map
	_robotList sync.Map
)

type Robot struct {
	CallCh      chan interface{}
	CloseCh     chan bool //机器人关闭
	WorkEndCh   chan bool //工作结束
	MatchTicker *time.Ticker
	GameInfo    map[string]interface{}
	User 		*model.User
}

type Call struct {
	GameId int
	Msg    interface{}
}

func New() *Robot {
	robot := &Robot{}
	callCh := make(chan interface{}, 1)
	closeCh := make(chan bool, 1)
	workEndCh := make(chan bool, 1)

	agent := &robotAgent{}
	agent.robotCallCh = callCh

	robot.User = &model.User{
		Uid:     robot.GetUniqueId(),
		Name:    robot.RandomName(),
		IsRobot: true,
	}
	robot.MatchTicker = time.NewTicker(conf.Get().Robot.MatchTime * time.Second)
	robot.CallCh = callCh
	robot.CloseCh = closeCh
	robot.WorkEndCh = workEndCh
	robot.User = robot.User.SetLoginInfo(robot.User, agent)
	_robotList.LoadOrStore(robot.User.Uid, robot.CloseCh)
	return robot
}

func (r *Robot) GetUniqueId() int {
	_robotId++
	return _robotId
}

func (r *Robot) RandomName() string {
	list := []string{
		"非凡哥", "鸡太美", "你好骚啊！", "小母牛飞上天", "跟非凡哥一桌",
		"撇嘴龙王", "一拳唐僧", "Dio", "JOJO", "坏老头", "今晚吃鸡",
	}
	return list[rand.Intn(len(list)-1)]
}



func (r *Robot) Close() {
	_robotList.Delete(r.User.Uid)
	r.User.DeleteCache(r.User.Uid)
}

//接收消息返回调用函数
func (r *Robot) Call(message interface{}) {
	callFunMap, ok := _callFun.Load(r.User.Game.Id)
	if !ok {
		return
	}

	callFunMap2 := callFunMap.(sync.Map)
	fun, ok := callFunMap2.Load(reflect.TypeOf(message))
	if !ok {
		return
	}

	fun.(func(interface{}, *Robot))(message, r)
}



func (r *Robot) CloseAll() {
	_robotList.Range(func(key, value interface{}) bool {
		if ch, ok := value.(chan bool); ok {
			ch <- true
		}
		return true
	})
}

//注册回调的消息调用的函数
func RegisterCallMsg(gameId int, message interface{}, fun func(interface{}, *Robot)) {
	callFunMap, ok := _callFun.Load(gameId)
	if !ok {
		callFunMap = sync.Map{}
	}
	callFunMap2 := callFunMap.(sync.Map)

	callFunMap2.Store(reflect.TypeOf(message), fun)
	_callFun.Store(gameId, callFunMap2)
}