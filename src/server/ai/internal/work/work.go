package work

import (
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/module"
	"reflect"
	"server/ai/internal/robot"
	"server/define"
	"server/game"
	"server/game/service/play"
	"server/gamedata"
	"server/models"
	"server/msg"
)

type Work struct {
	*robot.Robot
}

func init() {
	workMod := new(Work)
	for _, gameId := range play.AllGameId() {
		allGameHandler(workMod, gameId)
	}
}

//公共的非游戏模块
func commonHandler(w *Work) {

}

//所有游戏都注册的
func allGameHandler(work *Work, gameId int)  {
	work.RegisterCallMsg(gameId, &msg.S2C_MatchPlayer{}, work.MatchPlayerSuccess)
}

func (w *Work) Start(skeleton *module.Skeleton) {
	skeleton.Go( func() {
		work := new(Work)
		work.Robot = work.Create()
		for {
			select {
			case <-work.MatchTicker.C:
				work.MatchPlayer()
			case message := <-work.CallCh:
				work.Call(message)
			case <-work.WorkEndCh:
				log.Debug("机器人工作结束, 重新开启新的机器人")
				work.Restart(skeleton)
				return
			case <-work.CloseCh:
				log.Debug("机器人关闭,所有工作停止")
				work.Close()
				return
			}
		}
	}, nil)
}

func (w *Work) Restart(skeleton *module.Skeleton)  {
	w.Close()
	w.Start(skeleton)
}

func (w *Work) MatchPlayerSuccess(message interface{}, r *robot.Robot)  {
	r.MatchTicker.Stop()
}


func (w *Work) MatchPlayer()  {
	if w.Game.Status != define.GameFree {
		log.Debug("ROBOT已经匹配中或者游戏中")
		w.MatchTicker.Stop()
		return
	}

	//判断每个游戏的匹配人数然后加入
	allGameId := play.AllGameId()
	match := new(models.Match)
	for _, gameId := range allGameId{
		matchPlayerNum := gamedata.GetMatchNum(gameId)
		match2, _ := match.GameId2UidMap(gameId)
		//如果匹配缺少人加入
		if (matchPlayerNum - (int(match2.Num) % matchPlayerNum)) == 1 {
			send := &msg.C2S_MatchPlayer{GameId:gameId}
			game.ChanRPC.Go(reflect.TypeOf(send), send, *w.Agent)
			break
		}
	}

}