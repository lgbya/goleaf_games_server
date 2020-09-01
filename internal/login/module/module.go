package module

import (
	"server/internal/base"
	"server/internal/model"

	"github.com/name5566/leaf/module"
)

var (
	_skeleton = base.NewSkeleton()
	ChanRPC   = _skeleton.ChanRPCServer
)


type Module struct {
	*module.Skeleton
}

var _ module.Module = (*Module)(nil)
func (m *Module) OnInit() {
	m.Skeleton = _skeleton
}

func (m *Module) OnDestroy() {
	new(model.MaxUid).FlushDb()
}
