package mora

import (
	"server/internal/common/err-code"
	"server/internal/common/helper/game-helper"
	"server/internal/game/service/play"
	"server/internal/gate/protocol"
	"server/internal/model"

	"math"
)

type Mode struct {
	Info map[int]int
}

var _ play.Play = Mode{}
func (m Mode) Run(call *model.Call) {
	switch call.Msg.(type) {
	case *protocol.C2S_MoraPlay:
		m.handlePlay(call)
	}

}
func (m Mode) Start(room *model.Room, args ...interface{}) (map[string]interface{}, *model.Room) {
	room.GameInfo = Mode{Info: map[int]int{}}
	return make(map[string]interface{}), room
}

func (m Mode) Continue(user *model.User, room *model.Room, args ...interface{}) map[string]interface{} {
	continueInfo := make(map[string]interface{})
	userGameInfo := room.GameInfo.(Mode).Info[user.Uid]
	continueInfo["ply"] = userGameInfo
	return continueInfo
}

func (m Mode) handlePlay(call *model.Call) {

	//获取基本信息
	msg := call.Msg.(*protocol.C2S_MoraPlay)
	agent := call.Agent
	if !(msg.Ply == 1 || msg.Ply == 2 || msg.Ply == 3) {
		errCode.Msg(agent, "选择错误！")
		return
	}

	//修改角色缓存信息在游戏中
	user, room := gamehelper.CheckInRoom(agent)

	gameInfo := room.GameInfo.(Mode)
	gameInfo.Info[user.Uid] = msg.Ply
	room.GameInfo = gameInfo

	room.RoomId3Room(room)

	(*user.Agent).WriteMsg(&protocol.S2C_MoraPlay{
		Uid: user.Uid,
		Ply: msg.Ply,
	})
	//所有人都出完拳，判断输赢
	if len(gameInfo.Info) == len(room.UserList) {
		m.endGame(room)
	}
}

func (m Mode) endGame(room *model.Room) {
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

		(*user.Agent).WriteMsg(&protocol.S2C_EndGame{
			WinUid: winUid,
			End:    end,
		})
	}
	room.StopRoom()
}
