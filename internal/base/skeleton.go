package base

import (
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/module"
	"server/internal/common/conf"
)

func NewSkeleton() *module.Skeleton {
	config := conf.Get()
	skeleton := &module.Skeleton{
		GoLen:              config.Skeleton.GoLen,
		TimerDispatcherLen: config.Skeleton.TimerDispatcherLen,
		AsynCallLen:        config.Skeleton.AsynCallLen,
		ChanRPCServer:      chanrpc.NewServer(config.Skeleton.ChanRPCLen),
	}
	skeleton.Init()
	return skeleton
}
