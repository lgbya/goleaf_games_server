package service

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"server/game/internal"
	"server/game/service/common"
	"server/models"
	"server/msg"
)


func init() {
	internal.GetSkeleton().RegisterChanRPC("NewAgent", rpcNewAgent)
	internal.GetSkeleton().RegisterChanRPC("CloseAgent", rpcCloseAgent)
	internal.GetSkeleton().RegisterChanRPC("CommonInitGame", rpcCommonInitGame)

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
	}

	_ = agent
}

//
func rpcCommonInitGame(args []interface{}) {
	roomId := args[0].(int)
	gameId := args[1].(int)
	lUser := args[2].(map[int]*models.User)
	gameInfo := args[3]

	mUsers := make(map[int]msg.M_UserInfo)
	//告诉前端游戏开始
	for _, user := range lUser {
		mUsers[user.Uid] = msg.M_UserInfo{
			Uid : user.Uid,
			Name: user.Name,
		}
	}

	for _, user := range lUser {
		user.InRoomId = roomId
		user.Status = models.GamePlay
		user.Uid3User(user)
		(*user.Agent).WriteMsg(&msg.S2C_StartGame{
			RoomId :  roomId,
			UserList: mUsers,
		})
	}

	room := &models.Room{
		GameId:gameId, ID:roomId,
		UserList:lUser, GameInfo: gameInfo,
	}

	new(models.Room).RoomId3Room(roomId, room)
	if service, ok := NewGameService(gameId); ok {
		service.StartGame(room)
	}
	log.Debug("=======创建房间成功！！！房间号 %v==========", roomId)


}