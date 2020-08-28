package mora

import (
	"math"
	"server/game/service/common"
	"server/lib/tool/error2"
	"server/models"
	"server/msg"
)

type Mode struct {
	Info map[int]int
}


func (m *Mode) Run(call *models.Call) {
	switch call.Msg.(type)  {
	case *msg.C2S_MoraPlay:
		m.handlePlay(call)
	}

}
func (m *Mode) Start(room *models.Room, args ...interface{}) (map[string]interface{}, *models.Room) {
	room.GameInfo = Mode{Info: map[int]int{}}
	return make(map[string]interface{}), room
}

func (m *Mode) Continue(user *models.User, room *models.Room, args ...interface{}) map[string]interface{} {
	continueInfo := make(map[string]interface{})
	userGameInfo := room.GameInfo.(Mode).Info[user.Uid]
	continueInfo["ply"] = userGameInfo
	return continueInfo
}

func (m *Mode) handlePlay(call *models.Call){

	//获取基本信息
	message := call.Msg.(*msg.C2S_MoraPlay)
	agent := call.Agent
	if !(message.Ply == 1 || message.Ply == 2 || message.Ply == 3){
		error2.Msg(agent, "选择错误！")
		return
	}

	//修改角色缓存信息在游戏中
	user, room := common.CheckInRoom(agent)

	gameInfo := room.GameInfo.(Mode)
	gameInfo.Info[user.Uid] = message.Ply
	room.GameInfo = gameInfo

	room.RoomId3Room(room)

	(*user.Agent).WriteMsg(&msg.S2C_MoraPlay{
		Uid: user.Uid,
		Ply: message.Ply,
	})
	//所有人都出完拳，判断输赢
	if len(gameInfo.Info) == len(room.UserList) {
		m.endGame(room)
	}
}

func (m *Mode) endGame(room *models.Room)  {
	winUid, prePly := 0, 0
	gameInfo := make(map[int]int)
	for uid, ply := range room.GameInfo.(Mode).Info {
		gameInfo[uid] = ply
		if prePly != 0 {
			absPly := math.Abs(float64(ply - prePly))
			if (ply < prePly && absPly == 1) || (ply > prePly && absPly == 2) {
				winUid = uid
			} else if prePly == ply {
				winUid = 0
			}
		} else if prePly == 0 {
			winUid = uid
		}
		prePly = ply
	}

	end := make(map[string]interface{})
	end["gameInfo"] = gameInfo
	for _, user := range room.UserList {

		(*user.Agent).WriteMsg(&msg.S2C_EndGame{
			WinUid: winUid,
			End: end,
		})
	}
	room.StopRoom()
}