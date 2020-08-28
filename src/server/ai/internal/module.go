package internal

import (
	"github.com/name5566/leaf/module"
	"server/ai/internal/robot"
	"server/ai/internal/work"
	_ "server/ai/internal/work/more"
	_ "server/ai/internal/work/tictactoe"
	"server/base"
	"server/conf"
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
	workMod := new(work.Work)
	for i:=0; i< conf.RobotNum ;i++  {
		workMod.Start(skeleton)
	}

}

func (m *Module) Run(closeSig chan bool) {
	s := m.Skeleton
	newCloseSig := make(chan bool)
	s.Go(func() {
		select{
		case <-closeSig:
			//解散所有房间
			new(robot.Robot).CloseAll()
			newCloseSig<-true
		}
	}, nil)

	s.Run(newCloseSig)
}

func (m *Module) OnDestroy() {

}

func  GetSkeleton() *module.Skeleton {
	return skeleton
}
