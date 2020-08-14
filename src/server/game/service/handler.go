package service

import (
	"github.com/name5566/leaf/gate"
	"reflect"
	"server/game/internal"
	"server/game/service/common"
	"server/lib/tool/error2"
	"server/models"
	"server/msg"
)
//


func init() {
	handler(&msg.C2S_MatchPlayer{}, handleMatchPlayer)
	handler(&msg.C2S_CancelMatch{}, handleCancelMatch)
}


func handler(m interface{}, h interface{})  {
	internal.GetSkeleton().RegisterChanRPC(reflect.TypeOf(m), h)
}

func handleMatchPlayer(args []interface{}) {
	//获取基本信息
	protocol := args[0].(*msg.C2S_MatchPlayer)
	agent := args[1].(gate.Agent)

	//修改角色缓存信息在游戏中
	user, ok := common.CheckLogin(agent)
	if !ok{
		error2.Msg(agent,"请登录后再操作！")
		return
	}

	if user.Status != models.GameFree {
		error2.Msg(agent, "已经匹配或游戏中！")
		return
	}
	
	service, ok := NewGameService(protocol.GameId)
	if !ok  {
		error2.Msg(agent,  "游戏不存在！")
		return
	}
	user.Status = models.GameMath
	user.GameId = protocol.GameId
	user.Uid3User(user)

	service.MatchPlayer(user, protocol)

}

func handleCancelMatch(args []interface{})  {
	//获取基本信息
	protocol := args[0].(*msg.C2S_CancelMatch)
	agent := args[1].(gate.Agent)

	//修改角色缓存信息在游戏中
	user, ok := common.CheckLogin(agent)
	if ok == false || user.Status != models.GameMath{
		error2.Msg(agent, "未匹配游戏")
		return
	}
	service, ok := NewGameService(user.GameId)
	if !ok  {
		error2.Msg(agent,  "游戏不存在！")
		return
	}
	user.Status = models.GameFree
	user.GameId = 0
	user.Uid3User(user)
	service.CancelMatch(user, protocol)
}
