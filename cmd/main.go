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
	config := conf.Get()
	lconf.LogLevel = config.Log.Level
	lconf.LogPath = config.Log.Path
	lconf.LogFlag = config.Log.Flag
	lconf.ConsolePort = config.Server.ConsolePort
	lconf.ProfilePath = config.Server.ProfilePath

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
		ai.Module,
	)
}
