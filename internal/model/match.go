package model

import (
	"server/internal/common/cache"

	"sync"
	"sync/atomic"
)

type Match struct {
	GameId int
	Num    int64
	List   sync.Map
}

func (m *Match) ckGameId2UidMap(gameId int) string {
	return "key:game_id|value:uid_map/" + string(gameId)
}

func (m *Match) GameId2UidMap(gameId int) (*Match, bool) {
	if data, found := cache.New().Get(m.ckGameId2UidMap(gameId)); found {
		return data.(*Match), found
	}
	return m, false
}

func (m *Match) GameId3UidMap(gameId int, uid int) *Match {
	match, _ := m.GameId2UidMap(gameId)
	match.List.LoadOrStore(uid, uid)
	match.GameId = gameId
	atomic.AddInt64(&match.Num, 1)
	cache.New().SetNoExpiration(m.ckGameId2UidMap(gameId), match)
	return match
}

func (m *Match) GameId4UidMap(gameId int, uid int) {
	match, _ := m.GameId2UidMap(gameId)
	match.List.Delete(uid)
	match.GameId = gameId
	atomic.AddInt64(&match.Num, 1)
	cache.New().SetNoExpiration(m.ckGameId2UidMap(gameId), match)
}
