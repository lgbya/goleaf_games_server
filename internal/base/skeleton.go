package base

import (
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/module"
	"server/internal/common/conf"
)

func NewSkeleton() *module.Skeleton {
	skeleton := &module.Skeleton{
		GoLen:              conf.Server.Skeleton.GoLen,
		TimerDispatcherLen: conf.Server.Skeleton.TimerDispatcherLen,
		AsynCallLen:        conf.Server.Skeleton.AsynCallLen,
		ChanRPCServer:      chanrpc.NewServer(conf.Server.Skeleton.ChanRPCLen),
	}
	skeleton.Init()
	return skeleton
}
