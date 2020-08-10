package mora

import (
	"fmt"
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

var lMatch = make(map[int]int)

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

func (m *Mora) MatchPlayer(user *models.User, protocol interface{})  {
	gameId := protocol.(*msg.C2S_MatchPlayer).GameId

	//将当前角色uid加入对应的游戏匹配列表
	lMatch[user.Uid] = user.Uid

	//返回消息告诉前端已经加入匹配等待中
	(*user.Agent).WriteMsg(&msg.S2C_MatchPlayer{ GameId : gameId })

	//如果人数大于二人
	if len(lMatch) >= 2 {
		mapInt2User := make(map[int]*models.User)
		roomId := new(models.Room).GetUniqueID()
		modUser := new(models.User)
		for _, uid := range lMatch{
			if user, found	 := modUser.Uid2User(uid); found{
				delete(lMatch, user.Uid)
				mapInt2User[uid] = user
			}

			if len(mapInt2User) == 2 {
				log.Debug("============Start==========")
				more := Mora{Info: map[int]int{}}
				internal.ChanRPC.Go("CommonInitGame", roomId, gameId, mapInt2User, more)
				break
			}
		}

	}
}

func (m *Mora) CancelMatch(user *models.User, protocol interface{})  {
	delete(lMatch, user.Uid)
	(*user.Agent).WriteMsg(&msg.S2C_MatchPlayer{})
}

func (m *Mora) StartGame(room *models.Room)  {

}

func (m *Mora) handlerMoraPlaying(args []interface{}) {
	internal.GetSkeleton().Go(func() {
		//获取基本信息
		message := args[0].(*msg.C2S_MoraPlaying)
		agent := args[1].(gate.Agent)

		//修改角色缓存信息在游戏中
		user, found := common.CheckLogin(agent)
		if !found {
			error2.FatalMsg(agent, error2.ErrSystem, "请登录后再操作！")
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

		for _, user2 := range room.UserList {
			(*user2.Agent).WriteMsg(&msg.S2C_MoraPlaying{
				Uid: user.Uid,
				Ply: message.Ply,
			})
		}
		//所有人都出完拳，判断输赢
		if len(gameInfo.Info) == len(room.UserList) {
			m.endGame(room)
		}
	}, func() {})
}

func (m *Mora) endGame(room *models.Room)  {
	winUid, prePly := 0, 0
	for uid, ply := range room.GameInfo.(Mora).Info {
		if prePly != 0 {
			absPly := math.Abs(float64(ply - prePly))
			fmt.Println(absPly, ply, prePly)
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

	for _, user := range room.UserList {
		user, found := user.Uid2User(user.Uid)
		if found {
			user.InRoomId = 0
			user.Status = models.GameFree
			user.Uid3User(user)
		}

		(*user.Agent).WriteMsg(&msg.S2C_MoreReslut{
			WinUid: winUid,
		})
	}
	new(models.Room).RoomId4Room(room.ID)
}