package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"reflect"
	"server/game"
	"server/lib/tool/error2"
	"server/models"
	"server/msg"
	"time"
	"unicode/utf8"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&msg.C2S_Heart{}, handleHeart)
	handler(&msg.C2S_Register{}, handleRegister)
	handler(&msg.C2S_Login{}, handleLogin)
	handler(&msg.C2S_ResetLogin{}, handleResetLogin)
}

func handleHeart(args []interface{})  {
	agent := args[1].(gate.Agent)
	nowTime := time.Now().Unix()
	agent.WriteMsg(&msg.S2C_Heart{Time : nowTime})
}

func handleRegister(args []interface{})   {
	message := args[0].(*msg.C2S_Register)
	agent := args[1].(gate.Agent)

	if  utf8.RuneCountInString(message.Name) < 6 {
		error2.Msg(agent, "账号长度要6位或以上")
		return
	}

	if  utf8.RuneCountInString(message.Password) < 6 {
		error2.Msg(agent, "密码长度要6位或以上")
		return
	}

	if message.Password != message.ConfirmPassword {
		error2.Msg(agent, "确认密码错误")
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

	//注册成功，写入登录数据
	user = setLoginInfo(user, agent)
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

	user := new(models.User)
	user, ok := user.FindLoginName(message.Name)

	//防止重复登录和在不同设备登录
	if user.CheckRepeatLogin(user.Uid) {
		error2.Msg(agent, "你已经登录了！")
		return
	}
	
	if !ok {
		error2.Msg(agent, "该玩家名不存在！")
		return
	}

	if ! user.AuthLoginPassword(message.Password) {
		error2.Msg(agent, "登录密码错误！")
		return
	}

	if oldUser, found := user.Uid2User(user.Uid);found {
		user = oldUser
	}

	user = setLoginInfo(user, agent)
	agent.WriteMsg(&msg.S2C_Login{
		Uid: user.Uid,
		Name: user.Name,
		Gold: user.Gold,
		Token: user.Token,
		ExpiresTime: user.ExpiresAt,
	})

	game.ChanRPC.Go("ContinueGame", user)

}

func handleResetLogin(args []interface{}) {
	message := args[0].(*msg.C2S_ResetLogin)
	agent := args[1].(gate.Agent)


	user := new(models.User)
	user, ok := user.TempToken2User(message.Token)

	//如果断开
	if !ok {
		error2.FatalMsg(agent, error2.LoginInAgain, "链接超时！请重新登录")
		return
	}

	//防止重复登录和在不同设备登录
	if user.CheckRepeatLogin(user.Uid) {
		error2.Msg(agent, "你已经登录了！")
		return
	}

	if user.ExpiresAt <= time.Now().Unix() {
		error2.FatalMsg(agent, error2.LoginInAgain, "登录超时！请重新登录")
		return
	}

	user = setLoginInfo(user, agent)
	agent.WriteMsg(&msg.S2C_ResetLogin{
		Uid: user.Uid,
		Name: user.Name,
		Gold: user.Gold,
		Token: user.Token,
		ExpiresTime: user.ExpiresAt,
	})

	game.ChanRPC.Go("ContinueGame", user)

}

//写入登录数据
func setLoginInfo(user *models.User, agent gate.Agent)*models.User{
	user.GenerateToken()
	user.Agent = &agent
	agent.SetUserData(&models.Agent{ID:user.Uid, HeartTime : time.Now().Unix()})
	user.Uid3User(user)
	user.Common3LoginUid(user.Uid, &agent)
	return user
}
