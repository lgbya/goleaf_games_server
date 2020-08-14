package common

import (
	"github.com/name5566/leaf/gate"
	"server/models"
	"server/msg"
)

func CheckLogin(agent gate.Agent) (*models.User, bool) {
	if userAgent, ok := agent.UserData().(*models.Agent); ok {
		user, found := new(models.User).Uid2User(userAgent.ID)
		return user, found
	}

	return &models.User{}, false
}



func User2MUserList(lUser map[int]*models.User) map[int]msg.M_UserInfo {
	mUsers := make(map[int]msg.M_UserInfo)
	//告诉前端游戏开始
	for _, user := range lUser {
		mUsers[user.Uid] = msg.M_UserInfo{
			Uid : user.Uid,
			Name: user.Name,
		}
	}
	return  mUsers
}