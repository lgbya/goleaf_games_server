package game

import (
	"server/internal/game/module"
	_ "server/internal/game/service"
)

var (
	Module  = new(module.Module)
	ChanRPC = module.ChanRPC
)
