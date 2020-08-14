package mora

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"math"
	"reflect"
	"server/game/internal"
	"server/game/service/common"
	"server/lib/tool/error2"
	"server/models"
	"server/msg"
)



var matchList = make(map[int]int)

type Mora struct {
	Info map[int]int
}

func init() {
	m := new(Mora)
	handler(&msg.C2S_MoraPlaying{}, m.handlerMoraPlaying)
}

func handler(m interface{}, h interface{})  {
	internal.GetSkeleton().RegisterChanRPC(reflect.TypeOf(m), h)
}

func (m *Mora) MatchPlayer(user *models.User, args ...interface{})  {
	gameId := args[0].(*msg.C2S_MatchPlayer).GameId

	//将当前角色uid加入对应的游戏匹配列表
	matchList[user.Uid] = user.Uid

	//返回消息告诉前端已经加入匹配等待中
	(*user.Agent).WriteMsg(&msg.S2C_MatchPlayer{ GameId : gameId })

	//如果人数大于二人
	if len(matchList) >= 2 {
		userList := make(map[int]*models.User)
		roomId := new(models.Room).GetUniqueID()
		modUser := new(models.User)
		for _, uid := range matchList {
			if user, found	 := modUser.Uid2User(uid); found{
				delete(matchList, user.Uid)
				userList[uid] = user
			}

			if len(userList) == 2 {
				log.Debug("============Start==========")
				more := Mora{Info: map[int]int{}}
				internal.ChanRPC.Go("StartGame", roomId, gameId, userList, more)
				break
			}
		}

	}
}

func (m *Mora) CancelMatch(user *models.User, args ...interface{})  {
	delete(matchList, user.Uid)
	(*user.Agent).WriteMsg(&msg.S2C_CancelMatch{})
}

func (m *Mora) StartGame(room *models.Room, args ...interface{}) map[string]interface{} {
	return make(map[string]interface{})
}

func (m *Mora) ContinueGame(user *models.User, room *models.Room, args ...interface{}) map[string]interface{} {
	continueInfo := make(map[string]interface{})
	userGameInfo := room.GameInfo.(Mora).Info[user.Uid]
	continueInfo["ply"] = userGameInfo
	return continueInfo
}

func (m *Mora) handlerMoraPlaying(args []interface{}) {
	internal.GetSkeleton().Go(func() {


		//获取基本信息
		message := args[0].(*msg.C2S_MoraPlaying)
		agent := args[1].(gate.Agent)

		if !(message.Ply == 1 || message.Ply == 2 || message.Ply == 3){
			error2.Msg(agent, "选择错误！")
			return
		}
		
		//修改角色缓存信息在游戏中
		user, found := common.CheckLogin(agent)
		if !found {
			error2.Msg(agent, "请登录后再操作！")
			return
		}

		room, found := new(models.Room).RoomId2Room(user.InRoomId)
		if !found {
			error2.Msg(agent, "未加入游戏！")
			return
		}
		gameInfo := room.GameInfo.(Mora)
		gameInfo.Info[user.Uid] = message.Ply
		room.GameInfo = gameInfo

		room.RoomId3Room(user.InRoomId, room)

		(*user.Agent).WriteMsg(&msg.S2C_MoraPlaying{
			Uid: user.Uid,
			Ply: message.Ply,
		})
		//所有人都出完拳，判断输赢
		if len(gameInfo.Info) == len(room.UserList) {
			m.endGame(room)
		}
	}, func() {})
}

func (m *Mora) endGame(room *models.Room)  {
	winUid, prePly := 0, 0
	gameInfo := make(map[int]int)
	for uid, ply := range room.GameInfo.(Mora).Info {
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
		user, found := user.Uid2User(user.Uid)
		if found {
			user.InRoomId = 0
			user.Status = models.GameFree
			user.Uid3User(user)
		}

		(*user.Agent).WriteMsg(&msg.S2C_EndGame{
			WinUid: winUid,
			End: end,
		})
	}
	new(models.Room).RoomId4Room(room.ID)
}