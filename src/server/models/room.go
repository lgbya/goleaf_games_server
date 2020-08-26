package models

import (
	"server/lib/cache"
	"sync/atomic"
)

var roomId int64 = 100000

type Room struct {
	ID       int
	GameId   int
	UserList map[int]*User
	GameInfo interface{}
}

func (r *Room) GetUniqueID()  int {
	atomic.AddInt64(&roomId, 1)
	return int(roomId)
}

func (r *Room) ckRoomId2Room(roomId int) string {
	return "key:room_id|value:room/" + string(roomId)
}

func (r *Room) RoomId2Room(roomId int) (*Room, bool){
	if data, found := cache.New().Get(r.ckRoomId2Room(roomId));found{
		return data.(*Room), found
	}
	return nil, false
}

func (r *Room) RoomId3Room( data *Room){
	cache.New().SetNoExpiration(r.ckRoomId2Room(data.ID), data)
}

func (r *Room) RoomId4Room(roomId int ){
	cache.New().Delete(r.ckRoomId2Room(roomId))
}

func (r *Room) StopRoom()  {
	for _, user := range r.UserList {
		user, found := user.Uid2User(user.Uid)
		if found {
			user.InRoomId = 0
			user.Status = GameFree
			user.Uid3User(user)
		}
	}
	r.RoomId4Room(r.ID)
}