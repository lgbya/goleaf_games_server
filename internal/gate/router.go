package gate

import (
	"server/internal/game"
	"server/internal/login"
	"server/internal/protocol"
)

func init() {
	protocol.Processor.SetRouter(&protocol.C2S_Heart{}, login.ChanRPC)
	protocol.Processor.SetRouter(&protocol.C2S_Register{}, login.ChanRPC)
	protocol.Processor.SetRouter(&protocol.C2S_Login{}, login.ChanRPC)
	protocol.Processor.SetRouter(&protocol.C2S_ResetLogin{}, login.ChanRPC)
	protocol.Processor.SetRouter(&protocol.C2S_MatchPlayer{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.C2S_CancelMatch{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.C2S_MoraPlay{}, game.ChanRPC)
	protocol.Processor.SetRouter(&protocol.C2S_TictactoePlay{}, game.ChanRPC)
}
