package common

import (
	"github.com/name5566/leaf/gate"
	"server/models"
)

func CheckLogin(agent gate.Agent) (*models.User, bool) {
	if userAgent, ok := agent.UserData().(*models.Agent); ok {
		user, found := new(models.User).Uid2User(userAgent.ID)
		return user, found
	}

	return &models.User{}, false
}