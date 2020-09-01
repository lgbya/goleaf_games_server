package internal

import (
	"server/internal/ai/robot"
	"server/internal/ai/work"
	_ "server/internal/ai/work/more"
	_ "server/internal/ai/work/tictactoe"
	"server/internal/base"
	"server/internal/common/conf"

	"github.com/name5566/leaf/module"
)

var (
	_skeleton = base.NewSkeleton()
	ChanRPC   = _skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = _skeleton
	for i := 0; i < conf.Server.Robot.Num; i++ {
		work.Start(_skeleton)
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
