package service

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"server/define"
	"server/game/internal"
	"server/game/service/common"
	"server/game/service/play"
	"server/lib/tool/error2"
	"server/models"
	"server/msg"
	"time"
)

func init() {
	internal.GetSkeleton().RegisterChanRPC("NewAgent", rpcNewAgent)
	internal.GetSkeleton().RegisterChanRPC("CloseAgent", rpcCloseAgent)
	internal.GetSkeleton().RegisterChanRPC("StartGame", rpcStartGame)
	internal.GetSkeleton().RegisterChanRPC("ContinueGame", rpcContinueGame)
}

func rpcNewAgent(args []interface{}) {
	agent := args[0].(gate.Agent)
	_ = agent
}

func rpcCloseAgent(args []interface{}) {
	agent := args[0].(gate.Agent)
	if user, ok := common.CheckLogin(agent); ok{
		switch user.Status {
		case define.GameFree :
			user.DeleteCache(user.Uid)

		case define.GameMath:
			new(models.Match).GameId4UidMap(user.GameId, user.Uid)
			user.DeleteCache(user.Uid)
		}

		user.Common4LoginUid(user.Uid)
		user.TempToken3User(user, 30 * time.Second)
	}

	_ = agent
}

func rpcStartGame(args []interface{}) {

	internal.GetSkeleton().Go(func() {

		roomId := args[0].(int)
		gameId := args[1].(int)
		userList := args[2].(map[int]*models.User)

		callCh := make(chan models.Call, 10)
		stopCh := make(chan bool,1)

		playMod, _ := play.New(gameId)

		//修改房间信息
		room := &models.Room{
			GameId:gameId, ID:roomId,
			UserList: userList, CallCh: callCh,
			StopCh:stopCh,
		}

		//不同游戏的开始钩子
		startInfo, room := playMod.Start(room)
		new(models.Room).RoomId3Room(room)

		//通知房间内的所有玩家
		mUserList := common.User2MUserList(userList)
		for _, user := range userList {
			user.InRoomId = roomId
			user.Status = define.GamePlay
			user.Uid3User(user)
			(*user.Agent).WriteMsg(&msg.S2C_StartGame{
				RoomId :  roomId,
				GameId: gameId,
				UserList: mUserList,
				Start: startInfo,
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
	},nil)

}

func rpcContinueGame(args []interface{})  {
	user := args[0].(*models.User)
	if user.Status == define.GameFree {
		return
	}

	room, ok := new(models.Room).RoomId2Room(user.InRoomId)
	if  !ok {
		user.InRoomId = 0
		user.GameId = 0
		user.GameId = define.GameFree
		user.Uid3User(user)
		error2.Msg(*(user.Agent),  "游戏已经结束！")
		return
	}

	playMod, ok := play.New(user.GameId)
	if !ok  {
		error2.Msg(*(user.Agent),  "游戏不存在！")
		return
	}

	//不同游戏的开始继续游戏
	continueInfo := playMod.Continue(user, room)

	//通知重新进入游戏的玩家
	(*user.Agent).WriteMsg(&msg.S2C_ContinueGame{
		RoomId :  user.InRoomId,
		GameId : user.GameId,
		UserList: common.User2MUserList(room.UserList),
		Continue : continueInfo,
	})



}