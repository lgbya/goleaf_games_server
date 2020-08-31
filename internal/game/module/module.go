package module

import (
	"github.com/name5566/leaf/module"
	"server/internal/base"
	"server/internal/model"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func GetSkeleton() *module.Skeleton {
	return skeleton
}

func (m *Module) Run(closeSig chan bool) {
	s := m.Skeleton
	newCloseSig := make(chan bool)
	s.Go(func() {
		select {
		case <-closeSig:
			//解散所有房间
			new(model.Room).StopAllRoom()
			newCloseSig <- true
		}
	}, nil)

	s.Run(newCloseSig)
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
}

func (m *Module) OnDestroy() {

}
