package service

import (
	"github.com/name5566/leaf/gate"
	"reflect"
	"server/game/internal"
	"server/game/service/common"
	"server/gamedata"
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
	gameId := protocol.GameId

	//修改角色缓存信息在游戏中
	user, ok := common.CheckLogin(agent)
	if !ok{
		error2.FatalMsg(agent, error2.LoginInAgain, "请登录后再操作！")
		return
	}

	if user.Status != models.GameFree {
		error2.Msg(agent, "已经匹配或游戏中！")
		return
	}
	
	_, ok = NewGameService(gameId)
	if !ok  {
		error2.Msg(agent,  "游戏不存在！")
		return
	}

	matchPlayerNum := gamedata.GetMatchNum(gameId)
	if matchPlayerNum < 2{
		error2.Msg(agent,  "游戏配置错误！")
	}

	user.Status = models.GameMath
	user.GameId = gameId
	user.Uid3User(user)


	//将当前角色uid加入对应的游戏匹配列表
	match := new(models.Match)
	match = match.GameId3UidMap(gameId, user.Uid)

	//返回消息告诉前端已经加入匹配等待中
	(*user.Agent).WriteMsg(&msg.S2C_MatchPlayer{ GameId : gameId })



	if match.Num >= matchPlayerNum {
		userList := make(map[int]*models.User)
		roomId := new(models.Room).GetUniqueID()
		modUser := new(models.User)
		match.List.Range(func(uid, value interface{}) bool {
			if user, found	 := modUser.Uid2User(uid.(int)); found{
				userList[uid.(int)] = user
			}

			if len(userList) == matchPlayerNum {
				for uid := range userList{
					match.GameId4UidMap(match.GameId, uid)
				}
				internal.ChanRPC.Go("StartGame", roomId, match.GameId, userList)
				return false
			}
			return true
		})

	}



}

func handleCancelMatch(args []interface{})  {
	//获取基本信息
	agent := args[1].(gate.Agent)

	//修改角色缓存信息在游戏中
	user, ok := common.CheckLogin(agent)

	if !ok {
		error2.FatalMsg(agent, error2.LoginInAgain,"请登录后再操作！")
		return
	}

	if user.Status != models.GameMath{
		error2.Msg(agent, "未匹配游戏")
		return
	}

	user.Status = models.GameFree
	user.GameId = 0
	user.Uid3User(user)
	new(models.Match).GameId4UidMap(user.GameId, user.Uid)
}