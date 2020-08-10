package models

import (
	"server/lib/cache"
)

var roomId = 100000

type Room struct {
	ID       int
	GameId   int
	UserList map[int]*User
	GameInfo interface{}
}

func (r *Room) GetUniqueID()  int {
	return roomId + 1
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

func (r *Room) RoomId3Room(roomId int, data *Room){
	cache.New().SetNoExpiration(r.ckRoomId2Room(roomId), data)
}

func (r *Room) RoomId4Room(roomId int ){
	cache.New().Delete(r.ckRoomId2Room(roomId))
}

