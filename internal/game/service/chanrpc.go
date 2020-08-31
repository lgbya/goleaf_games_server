package service

import (
	"server/internal/common/define"
	"server/internal/common/err-code"
	"server/internal/common/helper/game-helper"
	"server/internal/game/module"
	"server/internal/game/service/play"
	"server/internal/model"
	"server/internal/protocol"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"

	"time"
)

func init() {
	module.GetSkeleton().RegisterChanRPC("NewAgent", rpcNewAgent)
	module.GetSkeleton().RegisterChanRPC("CloseAgent", rpcCloseAgent)
	module.GetSkeleton().RegisterChanRPC("StartGame", rpcStartGame)
	module.GetSkeleton().RegisterChanRPC("ContinueGame", rpcContinueGame)
}

func rpcNewAgent(args []interface{}) {
	agent := args[0].(gate.Agent)
	_ = agent
}

func rpcCloseAgent(args []interface{}) {
	agent := args[0].(gate.Agent)
	if user, ok := gamehelper.CheckLogin(agent); ok {
		switch user.Game.Status {
		case define.GameFree:
			user.DeleteCache(user.Uid)

		case define.GameMath:
			new(model.Match).GameId4UidMap(user.Game.Id, user.Uid)
			user.DeleteCache(user.Uid)
		}

		user.Common4LoginUid(user.Uid)
		user.TempToken3User(user, 30*time.Second)
	}

	_ = agent
}

func rpcStartGame(args []interface{}) {

	module.GetSkeleton().Go(func() {

		roomId := args[0].(int)
		gameId := args[1].(int)
		userList := args[2].(map[int]*model.User)

		callCh := make(chan model.Call, 10)
		stopCh := make(chan bool, 1)

		playMod, _ := play.New(gameId)

		//修改房间信息
		room := &model.Room{
			GameId: gameId, ID: roomId,
			UserList: userList, CallCh: callCh,
			StopCh: stopCh,
		}

		//不同游戏的开始钩子
		startInfo, room := playMod.Start(room)
		new(model.Room).RoomId3Room(room)

		//通知房间内的所有玩家
		mUserList := gamehelper.User2MUserList(userList)
		for _, user := range userList {
			user.Game.InRoomId = roomId
			user.Game.Status = define.GamePlay
			user.Uid3User(user)
			(*user.Agent).WriteMsg(&protocol.S2C_StartGame{
				RoomId:   roomId,
				GameId:   gameId,
				UserList: mUserList,
				Start:    startInfo,
			})
		}

		log.Debug("=======创建房间成功！！！房间号 %v==========", roomId)
		isStopRoom := false
		for !isStopRoom {
			select {
			case call := <-callCh:
				playMod.Run(&call)
			case <-stopCh:
				log.Debug("========房间解散%v========", roomId)
				isStopRoom = true
			}
		}
	}, nil)

}

func rpcContinueGame(args []interface{}) {
	user := args[0].(*model.User)
	if user.Game.Status == define.GameFree {
		return
	}

	room, ok := new(model.Room).RoomId2Room(user.Game.InRoomId)
	if !ok {
		user.Game.InRoomId = 0
		user.Game.Id = 0
		user.Game.Status = define.GameFree
		user.Uid3User(user)
		errCode.Msg(*(user.Agent), "游戏已经结束！")
		return
	}

	playMod, ok := play.New(user.Game.Id)
	if !ok {
		errCode.Msg(*(user.Agent), "游戏不存在！")
		return
	}

	//不同游戏的开始继续游戏
	continueInfo := playMod.Continue(user, room)

	//通知重新进入游戏的玩家
	(*user.Agent).WriteMsg(&protocol.S2C_ContinueGame{
		RoomId:   user.Game.InRoomId,
		GameId:   user.Game.Id,
		UserList: gamehelper.User2MUserList(room.UserList),
		Continue: continueInfo,
	})

}
