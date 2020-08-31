package internal

import (
	"server/internal/ai/module/robot"
	"server/internal/ai/module/work"
	_ "server/internal/ai/module/work/more"
	_ "server/internal/ai/module/work/tictactoe"
	"server/internal/base"
	"server/internal/common/conf"

	"github.com/name5566/leaf/module"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	workMod := work.Work{}
	for i := 0; i < conf.Server.Robot.Num; i++ {
		workMod.Start(skeleton)
	}

}

func (m *Module) Run(closeSig chan bool) {
	s := m.Skeleton
	newCloseSig := make(chan bool)
	s.Go(func() {
		select {
		case <-closeSig:
			//解散所有房间
			new(robot.Robot).CloseAll()
			newCloseSig <- true
		}
	}, nil)

	s.Run(newCloseSig)
}

func (m *Module) OnDestroy() {

}
