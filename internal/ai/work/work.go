package work

import (
	"server/internal/ai/robot"
	"server/internal/common/define"
	"server/internal/common/gamedata"
	"server/internal/game"
	"server/internal/game/service"
	"server/internal/model"
	"server/internal/protocol"

	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/module"

	"reflect"
	"sync"
)

var _lock sync.Mutex

type Work struct {
	*robot.Robot
}

func init() {
	workMod := Work{}
	for _, gameId := range service.AllGameId() {
		allGameHandler(workMod, gameId)
	}
}

//公共的非游戏模块
func commonHandler(w *Work) {

}

//所有游戏都注册的
func allGameHandler(work Work, gameId int) {
	robot.RegisterCallMsg(gameId, &protocol.S2C_MatchPlayer{}, work.MatchPlayerSuccess)
}

func Start(skeleton *module.Skeleton) {
	skeleton.Go(func() {
		w := &Work{}
		w.Robot = robot.New()
		for {
			select {
			case <-w.MatchTicker.C:
				w.MatchPlayer()
			case message := <-w.CallCh:
				w.Call(message)
			case <-w.WorkEndCh:
				log.Debug("机器人工作结束, 重新开启新的机器人")
				w.Restart(skeleton)
				return
			case <-w.CloseCh:
				//log.Debug("机器人关闭,所有工作停止")
				w.Close()
				return
			}
		}
	}, nil)
}

func (w Work) Restart(skeleton *module.Skeleton) {
	w.Close()
	Start(skeleton)
}

func (w Work) MatchPlayerSuccess(message interface{}, r *robot.Robot) {
	r.MatchTicker.Stop()
}

func (w Work) MatchPlayer() {

	if w.User.Game.Status != define.GameFree {
		log.Debug("ROBOT已经匹配中或者游戏中")
		w.MatchTicker.Stop()
		return
	}

	defer _lock.Unlock()
	_lock.Lock()
	//判断每个游戏的匹配人数然后加入
	allGameId := service.AllGameId()
	match := model.Match{}
	for _, gameId := range allGameId {
		matchPlayerNum := gamedata.GetMatchNum(gameId)

		if match2, ok := match.GameId2UidMap(gameId);ok{

			//如果匹配缺少人加入
			if (matchPlayerNum - (int(match2.Num) % matchPlayerNum)) == 1 {
				send := &protocol.C2S_MatchPlayer{GameId: gameId}
				game.ChanRPC.Go(reflect.TypeOf(send), send, *w.User.Agent)
				return
			}
		}
	}

}
