package module

import (
	"server/internal/common/define"
	"server/internal/common/err-code"
	"server/internal/common/gamedata"
	"server/internal/common/helper/game-helper"
	"server/internal/game/service"
	"server/internal/gate/protocol"
	"server/internal/model"

	"github.com/name5566/leaf/gate"

	"reflect"
)

//

func init() {
	handler(&protocol.C2S_MatchPlayer{}, handleMatchPlayer)
	handler(&protocol.C2S_CancelMatch{}, handleCancelMatch)
	handler(&protocol.C2S_TictactoePlay{}, handlePlay)
	handler(&protocol.C2S_MoraPlay{}, handlePlay)
}

func handler(m interface{}, h interface{}) {
	_skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func handlePlay(args []interface{}) {
	_skeleton.Go(func() {

		//获取基本信息
		message := args[0]
		agent := args[1].(gate.Agent)

		//修改角色缓存信息在游戏中
		user, room := gamehelper.CheckInRoom(agent)
		if user == nil || room == nil {
			errCode.FatalMsg(agent, errCode.LoginInAgain, "未加入游戏！")
			return
		}
		ch := room.CallCh
		ch <- model.Call{Uid: user.Uid, Agent: agent, Msg: message}

	}, nil)
}

func handleMatchPlayer(args []interface{}) {
	//获取基本信息
	msg := args[0].(*protocol.C2S_MatchPlayer)
	agent := args[1].(gate.Agent)
	gameId := msg.GameId

	//修改角色缓存信息在游戏中
	user, ok := gamehelper.CheckLogin(agent)
	if !ok {
		errCode.FatalMsg(agent, errCode.LoginInAgain, "请登录后再操作！")
		return
	}

	if user.Game.Status != define.GameFree {
		errCode.Msg(agent, "已经匹配或游戏中！")
		return
	}

	if _, ok = service.NewPlay(gameId); !ok {
		errCode.Msg(agent, "游戏不存在！")
		return
	}

	matchPlayerNum := gamedata.GetMatchNum(gameId)
	if matchPlayerNum < 2 {
		errCode.Msg(agent, "游戏配置错误！")
	}

	user.Game.Status = define.GameMath
	user.Game.Id = gameId
	user.Uid3User(user)

	//将当前角色uid加入对应的游戏匹配列表
	match := &model.Match{}
	match = match.GameId3UidMap(gameId, user.Uid)

	//返回消息告诉前端已经加入匹配等待中
	(*user.Agent).WriteMsg(&protocol.S2C_MatchPlayer{GameId: gameId})

	if int(match.Num) >= matchPlayerNum {
		userList := make(map[int]*model.User)
		roomId := new(model.Room).GetUniqueID()
		modUser := new(model.User)
		match.List.Range(func(uid, value interface{}) bool {
			if user, found := modUser.Uid2User(uid.(int)); found {
				userList[uid.(int)] = user
			}

			if len(userList) == matchPlayerNum {
				for uid := range userList {
					match.GameId4UidMap(match.GameId, uid)
				}
				ChanRPC.Go("StartGame", roomId, match.GameId, userList)
				return false
			}
			return true
		})

	}

}

func handleCancelMatch(args []interface{}) {
	//获取基本信息
	agent := args[1].(gate.Agent)

	//修改角色缓存信息在游戏中
	user, ok := gamehelper.CheckLogin(agent)

	if !ok {
		errCode.FatalMsg(agent, errCode.LoginInAgain, "请登录后再操作！")
		return
	}

	if user.Game.Status != define.GameMath {
		errCode.Msg(agent, "未匹配游戏")
		return
	}

	user.Game.Id = 0
	user.Game.Status = define.GameFree
	user.Uid3User(user)

	match := model.Match{}
	match.GameId4UidMap(user.Game.Id, user.Uid)
}
