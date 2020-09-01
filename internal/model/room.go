package model

import (
	"server/internal/common/cache"
	"server/internal/common/define"

	"github.com/name5566/leaf/gate"

	"sync"
	"sync/atomic"
)

var _roomId int64 = 100000
var roomIdMap sync.Map

type Room struct {
	ID       int
	GameId   int
	UserList map[int]*User
	GameInfo interface{}
	CallCh   chan Call
	StopCh   chan bool
}

type Call struct {
	Uid   int
	Agent gate.Agent
	Msg   interface{}
}

func (r Room) GetUniqueID() int {
	atomic.AddInt64(&_roomId, 1)
	return int(_roomId)
}

func (r Room) ckRoomId2Room(roomId int) string {
	return "key:room_id|value:room/" + string(roomId)
}

func (r Room) RoomId2Room(roomId int) (*Room, bool) {
	if data, found := cache.New().Get(r.ckRoomId2Room(roomId)); found {
		return data.(*Room), found
	}
	return nil, false
}

func (r Room) RoomId3Room(data *Room) {
	roomIdMap.LoadOrStore(data.ID, data.StopCh)
	cache.New().SetNoExpiration(r.ckRoomId2Room(data.ID), data)
}

func (r Room) RoomId4Room(roomId int) {
	roomIdMap.Delete(roomId)
	cache.New().Delete(r.ckRoomId2Room(roomId))
}

func (r Room) StopAllRoom() {
	roomIdMap.Range(func(key, value interface{}) bool {
		if ch, ok := value.(chan bool); ok {
			ch <- true
		}
		return true
	})
}

func (r Room) StopRoom() {
	for _, user := range r.UserList {
		user, found := user.Uid2User(user.Uid)
		if found {
			user.Game.InRoomId = 0
			user.Game.Status = define.GameFree
			user.Uid3User(user)
		}
	}

	r.StopCh <- true

	close(r.CallCh)
	close(r.StopCh)
	r.RoomId4Room(r.ID)
}
