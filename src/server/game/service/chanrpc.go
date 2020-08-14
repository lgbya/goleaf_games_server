package service

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"server/game/internal"
	"server/game/service/common"
	"server/lib/tool/error2"
	"server/models"
	"server/msg"
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
		case models.GameFree :
			user.DeleteCache(user.Uid)

		case models.GameMath:
			if service, ok := NewGameService(user.GameId);ok  {
				service.CancelMatch(user, &msg.S2C_CancelMatch{})
			}
			user.DeleteCache(user.Uid)
		}

		user.Common4LoginUid(user.Uid)
	}

	_ = agent
}

//
func rpcStartGame(args []interface{}) {
	roomId := args[0].(int)
	gameId := args[1].(int)
	userList := args[2].(map[int]*models.User)
	gameInfo := args[3]

	//修改房间信息
	room := &models.Room{
		GameId:gameId, ID:roomId,
		UserList: userList, GameInfo: gameInfo,
	}
	new(models.Room).RoomId3Room(roomId, room)
	service, ok := NewGameService(gameId)
	if !ok  {
		for _, user := range userList {
			user.InRoomId = 0
			user.Status = models.GameFree
			user.GameId = 0
			user.Uid3User(user)
			error2.Msg(*(user.Agent),  "游戏不存在！")
		}
		return
	}
	//不同游戏的开始钩子
	startInfo := service.StartGame(room)

	//通知房间内的所有玩家
	mUserList := common.User2MUserList(userList)
	for _, user := range userList {
		user.InRoomId = roomId
		user.Status = models.GamePlay
		user.Uid3User(user)
		(*user.Agent).WriteMsg(&msg.S2C_StartGame{
			RoomId :  roomId,
			UserList: mUserList,
			Start: startInfo,
		})
	}

	log.Debug("=======创建房间成功！！！房间号 %v==========", roomId)

}

func rpcContinueGame(args []interface{})  {
	user := args[0].(*models.User)
	if user.Status == models.GameFree {
		return
	}

	room, ok := new(models.Room).RoomId2Room(user.InRoomId)
	if  !ok {
		user.InRoomId = 0
		user.GameId = 0
		user.GameId = models.GameFree
		user.Uid3User(user)
		error2.Msg(*(user.Agent),  "游戏已经结束！")
		return
	}

	service, ok := NewGameService(user.GameId)
	if !ok  {
		error2.Msg(*(user.Agent),  "游戏不存在！")
		return
	}

	//不同游戏的开始继续游戏
	continueInfo := service.ContinueGame(user, room)

	//通知重新进入游戏的玩家
	(*user.Agent).WriteMsg(&msg.S2C_ContinueGame{
		RoomId :  user.InRoomId,
		GameId : user.GameId,
		UserList: common.User2MUserList(room.UserList),
		Continue : continueInfo,
	})



}