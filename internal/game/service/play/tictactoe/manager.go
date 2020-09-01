package tictactoe

import (
	"server/internal/common/err-code"
	"server/internal/common/helper/game-helper"
	"server/internal/game/service/play"
	"server/internal/model"
	"server/internal/protocol"

	"github.com/chenhg5/collection"
	"github.com/shopspring/decimal"
)

type Mode struct {
	Info       map[int]player `json:"info"`
	Blank      []int          `json:"blank"`
	CurrentUid int            `json:"currentUid"`
}

type player struct {
	Icon int   `json:"icon"`
	List []int `json:"list"`
}

var _ play.Play = Mode{}
func (m Mode) Run(call *model.Call) {
	switch call.Msg.(type) {
	case *protocol.C2S_TictactoePlay:
		m.handlePlay(call)
	}
}

func (m Mode) Start(room *model.Room, args ...interface{}) (map[string]interface{}, *model.Room) {
	//为玩家分配对号和叉号
	gameInfo := Mode{Info: map[int]player{}, Blank: []int{1, 2, 3, 4, 5, 6, 7, 8, 9}}
	icon := 0
	userIcon := make(map[int]int)
	for uid := range room.UserList {
		userIcon[uid] = icon
		gameInfo.Info[uid] = player{
			Icon: icon,
			List: []int{},
		}
		icon++
		gameInfo.CurrentUid = uid
	}
	room.GameInfo = gameInfo
	startInfo := map[string]interface{}{
		"currentUid": gameInfo.CurrentUid,
		"userIcon":   userIcon,
	}
	return startInfo, room
}

func (m Mode) Continue(user *model.User, room *model.Room, args ...interface{}) map[string]interface{} {
	continueInfo := make(map[string]interface{})
	gameInfo := room.GameInfo.(Mode)
	continueInfo["mode"] = gameInfo
	return continueInfo
}

func (m Mode) handlePlay(call *model.Call) {

	//获取基本信息
	message := call.Msg.(*protocol.C2S_TictactoePlay)
	agent := call.Agent

	//修改角色缓存信息在游戏中
	user, room := gamehelper.CheckInRoom(agent)

	number := message.Number
	gameInfo := room.GameInfo.(Mode)

	result := collection.Collect(gameInfo.Blank).Contains(number)
	if !(number > 0 && result) {
		errCode.Msg(agent, "选择错误！")
		return
	}

	if gameInfo.CurrentUid != user.Uid {
		errCode.Msg(agent, "未轮到你操作！")
		return
	}

	//记录当前玩家选择的号码
	playerInfo := gameInfo.Info[user.Uid]
	playerInfo.List = append(playerInfo.List, number)
	gameInfo.Info[user.Uid] = playerInfo
	for uid := range gameInfo.Info {
		if gameInfo.CurrentUid != uid {
			gameInfo.CurrentUid = uid
			break
		}
	}

	//请求已经选择的号码
	gameInfo.Blank = collection.Collect(gameInfo.Blank).Reject(func(item, value interface{}) bool {
		return value.(decimal.Decimal).IntPart() == int64(number)
	}).ToIntArray()

	//修改房间信息
	room.GameInfo = gameInfo
	room.RoomId3Room(room)

	//通知前端修改
	for _, accUser := range room.UserList {
		(*accUser.Agent).WriteMsg(&protocol.S2C_TictactoePlay{
			Uid:        user.Uid,
			Number:     number,
			CurrentUid: gameInfo.CurrentUid,
		})
	}

	//判断胜负
	if result, winUid, winCombine := m.checkVictory(user.Uid, gameInfo); result {
		m.endGame(winUid, winCombine, room)
	}

}

func (m Mode) checkVictory(uid int, gameInfo Mode) (bool, int, []int) {

	playerInfo := gameInfo.Info[uid]

	combineList := [...][]int{
		{1, 2, 3}, {4, 5, 6}, {7, 8, 9},
		{1, 4, 7}, {2, 5, 8}, {3, 6, 9},
		{1, 5, 9}, {3, 5, 7},
	}

	for _, combine := range combineList {
		result := collection.Collect(combine).Every(func(item, value interface{}) bool {
			value = int(value.(decimal.Decimal).IntPart())
			return collection.Collect(playerInfo.List).Contains(value)
		})

		if result {
			return true, uid, combine
		}
	}

	if len(gameInfo.Blank) == 0 {
		return true, 0, []int{}
	}

	return false, 0, []int{}

}

func (m Mode) endGame(winUid int, winCombine []int, room *model.Room) {

	end := map[string]interface{}{
		"winUid":     winUid,
		"winCombine": winCombine,
	}

	for _, user := range room.UserList {
		(*user.Agent).WriteMsg(&protocol.S2C_EndGame{
			WinUid: winUid,
			End:    end,
		})
	}
	room.StopRoom()
}
