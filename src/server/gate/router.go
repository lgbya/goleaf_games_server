package gate

import (
	"server/game"
	"server/login"
	"server/msg"
)

func init() {
	msg.Processor.SetRouter(&msg.C2S_Register{}, login.ChanRPC)
	msg.Processor.SetRouter(&msg.C2S_Login{}, login.ChanRPC)
	msg.Processor.SetRouter(&msg.C2S_MatchPlayer{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.C2S_MoraPlaying{}, game.ChanRPC)


}