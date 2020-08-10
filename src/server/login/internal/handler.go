package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"reflect"
	"server/lib/tool/error2"
	"server/models"
	"server/msg"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&msg.C2S_Register{}, handleRegister)
	handler(&msg.C2S_Login{}, handleLogin)
}

func handleRegister(args []interface{})   {
	message := args[0].(*msg.C2S_Register)
	agent := args[1].(gate.Agent)

	if message.Password != message.ConfirmPassword {
		error2.Msg(agent, "")
		return
	}

	user := new(models.User)
	if _, ok := user.FindLoginName(message.Name); ok {
		error2.Msg(agent, "该玩家名已经存在！")
		return
	}

	user, err := user.Create(message.Name, message.Password)
	if err != nil{
		log.Debug("注册失败原因:%v", err)
		error2.Msg(agent, "注册失败！")
		return
	}

	//注册成功，写入数据
	agent.SetUserData(&models.Agent{ID:user.Uid})
	user.Agent = &agent
	user.Uid3User(user)

	//返回注册成功消息
	agent.WriteMsg(&msg.S2C_Register{
		Uid: user.Uid,
		Name: user.Name,
		Gold: user.Gold,
		Token: user.Token,
		ExpiresTime: user.ExpiresAt,
	})

}

func handleLogin(args []interface{})  {
	message := args[0].(*msg.C2S_Login)
	agent := args[1].(gate.Agent)

	if agent.UserData() != nil {
		error2.Msg(agent, "你已经登录了！")
		return
	}

	user := new(models.User)
	user, ok := user.FindLoginName(message.Name)

	if !ok {
		error2.Msg(agent, "该玩家名不存在！")
		return
	}

	if user.AuthLoginPassword(message.Password) {
		error2.Msg(agent, "登录密码错误！")
		return
	}

	oldUser, found := user.Uid2User(user.Uid)
	if found  {
		if oldUser.Status == models.GamePlay {
			user = oldUser
		}else{
			error2.Msg(agent, "该玩家在其他设备已登录！")
			return
		}
	}

	user.Agent = &agent
	agent.SetUserData(&models.Agent{ID:user.Uid})
	user.Uid3User(user)

	agent.WriteMsg(&msg.S2C_Login{
		Uid: user.Uid,
		Name: user.Name,
		Gold: user.Gold,
		Token: user.Token,
		ExpiresTime: user.ExpiresAt,
	})

}