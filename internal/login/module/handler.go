package module

import (
	"server/internal/common/err-code"
	"server/internal/game"
	"server/internal/gate/protocol"
	"server/internal/model"

	"github.com/name5566/leaf/gate"

	"reflect"
	"time"
	"unicode/utf8"
)

func handler(m interface{}, h interface{}) {
	_skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&protocol.C2S_Heart{}, handleHeart)
	handler(&protocol.C2S_Register{}, handleRegister)
	handler(&protocol.C2S_Login{}, handleLogin)
	handler(&protocol.C2S_ResetLogin{}, handleResetLogin)
}

func handleHeart(args []interface{}) {
	agent := args[1].(gate.Agent)
	nowTime := time.Now().Unix()
	agent.WriteMsg(&protocol.S2C_Heart{Time: nowTime})
}

func handleRegister(args []interface{}) {
	message := args[0].(*protocol.C2S_Register)
	agent := args[1].(gate.Agent)

	if utf8.RuneCountInString(message.Name) < 6 {
		errCode.Msg(agent, "账号长度要6位或以上")
		return
	}

	if utf8.RuneCountInString(message.Password) < 6 {
		errCode.Msg(agent, "密码长度要6位或以上")
		return
	}

	if message.Password != message.ConfirmPassword {
		errCode.Msg(agent, "确认密码错误")
		return
	}

	user := &model.User{}
	if _, ok := user.FindLoginName(message.Name); ok {
		errCode.Msg(agent, "该玩家名已经存在！")
		return
	}

	if user.Create(message.Name, message.Password) != nil {
		errCode.Msg(agent, "注册失败！")
		return
	}

	//注册成功，写入登录数据
	user = user.SetLoginInfo(user, agent)
	//返回注册成功消息
	agent.WriteMsg(&protocol.S2C_Register{
		Uid:         user.Uid,
		Name:        user.Name,
		Gold:        user.Gold,
		Token:       user.Token,
		ExpiresTime: user.ExpiresAt,
	})

}

func handleLogin(args []interface{}) {
	message := args[0].(*protocol.C2S_Login)
	agent := args[1].(gate.Agent)

	user := &model.User{}
	user, ok := user.FindLoginName(message.Name)
	if !ok {
		errCode.Msg(agent, "该玩家名不存在！")
		return
	}

	//防止重复登录和在不同设备登录
	if user.CheckRepeatLogin(user.Uid) {
		errCode.Msg(agent, "你已经登录了！")
		return
	}

	if !user.AuthLoginPassword(message.Password) {
		errCode.Msg(agent, "登录密码错误！")
		return
	}

	if oldUser, found := user.Uid2User(user.Uid); found {
		user = oldUser
	}

	user = user.SetLoginInfo(user, agent)
	agent.WriteMsg(&protocol.S2C_Login{
		Uid:         user.Uid,
		Name:        user.Name,
		Gold:        user.Gold,
		Token:       user.Token,
		ExpiresTime: user.ExpiresAt,
	})

	game.ChanRPC.Go("ContinueGame", user)

}

func handleResetLogin(args []interface{}) {
	message := args[0].(*protocol.C2S_ResetLogin)
	agent := args[1].(gate.Agent)

	user := &model.User{}
	user, ok := user.TempToken2User(message.Token)

	//如果断开
	if !ok {
		errCode.FatalMsg(agent, errCode.LoginInAgain, "链接超时！请重新登录")
		return
	}

	//防止重复登录和在不同设备登录
	if user.CheckRepeatLogin(user.Uid) {
		errCode.Msg(agent, "你已经登录了！")
		return
	}

	if user.ExpiresAt <= time.Now().Unix() {
		errCode.FatalMsg(agent, errCode.LoginInAgain, "登录超时！请重新登录")
		return
	}

	user = user.SetLoginInfo(user, agent)
	agent.WriteMsg(&protocol.S2C_ResetLogin{
		Uid:         user.Uid,
		Name:        user.Name,
		Gold:        user.Gold,
		Token:       user.Token,
		ExpiresTime: user.ExpiresAt,
	})

	game.ChanRPC.Go("ContinueGame", user)

}
