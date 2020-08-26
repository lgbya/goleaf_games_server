package models

import (
	"server/lib/cache"
	"sync"
)

var matchLock sync.RWMutex

type Match struct {
	GameId int
	Num int
	List sync.Map
}

func (m *Match) ckGameId2UidMap(gameId int) string {
	return "key:game_id|value:uid_map/" + string(gameId)
}

func (m *Match) GameId2UidMap(gameId int) (*Match, bool) {
	if data, found := cache.New().Get(m.ckGameId2UidMap(gameId)); found{
		return data.(*Match), found
	}
	return m, false
}

func (m *Match) GameId3UidMap(gameId int,  uid int) *Match {
	match, _ := m.GameId2UidMap(gameId)
	match.List.LoadOrStore(uid, uid)
	matchLock.Lock()
	match.GameId = gameId
	match.Num++
	matchLock.Unlock()
	cache.New().SetNoExpiration(m.ckGameId2UidMap(gameId), match)
	return match
}

func (m *Match) GameId4UidMap(gameId int,  uid int ){
	match, _ := m.GameId2UidMap(gameId)
	match.List.Delete(uid)
	matchLock.Lock()
	match.GameId = gameId
	match.Num++
	matchLock.Unlock()
	cache.New().SetNoExpiration(m.ckGameId2UidMap(gameId), match)
}
