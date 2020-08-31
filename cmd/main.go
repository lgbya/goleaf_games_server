package main

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"server/internal/ai"
	"server/internal/common/conf"
	"server/internal/game"
	"server/internal/gate"
	"server/internal/login"
)

func main() {

	//初始化基础配置
	lconf.LogLevel = conf.Server.Log.Level
	lconf.LogPath = conf.Server.Log.Path
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
		ai.Module,
	)
}
