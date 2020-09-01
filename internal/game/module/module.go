package module

import (
	"github.com/name5566/leaf/module"
	"server/internal/base"
	"server/internal/model"
)

var (
	_skeleton = base.NewSkeleton()
	ChanRPC   = _skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

var _ module.Module = (*Module)(nil)
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
	m.Skeleton = _skeleton
}

func (m *Module) OnDestroy() {

}

func GetSkeleton() *module.Skeleton {
	return _skeleton
}