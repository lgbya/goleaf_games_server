package gamehelper

import (
	"github.com/name5566/leaf/gate"
	"server/internal/model"
	"server/internal/protocol"
)

func CheckLogin(agent gate.Agent) (*model.User, bool) {
	if userAgent, ok := agent.UserData().(*model.Agent); ok {
		user, found := new(model.User).Uid2User(userAgent.ID)
		return user, found
	}

	return &model.User{}, false
}

func User2MUserList(lUser map[int]*model.User) map[int]protocol.M_UserInfo {
	mUsers := make(map[int]protocol.M_UserInfo)
	//告诉前端游戏开始
	for _, user := range lUser {
		mUsers[user.Uid] = protocol.M_UserInfo{
			Uid:  user.Uid,
			Name: user.Name,
		}
	}
	return mUsers
}

func CheckInRoom(agent gate.Agent) (*model.User, *model.Room) {

	userAgent, ok := agent.UserData().(*model.Agent)
	if !ok {
		return nil, nil
	}

	user, ok := new(model.User).Uid2User(userAgent.ID)
	if !ok {
		return nil, nil
	}

	room, ok := new(model.Room).RoomId2Room(user.Game.InRoomId)
	if !ok {
		return nil, nil
	}

	return user, room
}
