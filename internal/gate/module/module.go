package module

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/module"
	"server/internal/common/conf"
	"server/internal/game"
	"server/internal/protocol"
	"time"
)

type Module struct {
	*gate.Gate
}

var _ module.Module = (*Module)(nil)
func (m *Module) OnInit() {
	m.Gate = &gate.Gate{
		MaxConnNum:      conf.Server.MaxConnNum,
		PendingWriteNum: conf.Server.Gate.PendingWriteNum,
		MaxMsgLen:       conf.Server.Gate.MaxMsgLen,
		WSAddr:          conf.Server.WSAddr,
		HTTPTimeout:     conf.Server.Gate.HTTPTimeout * time.Second,
		CertFile:        conf.Server.CertFile,
		KeyFile:         conf.Server.KeyFile,
		TCPAddr:         conf.Server.TCPAddr,
		LenMsgLen:       conf.Server.Gate.LenMsgLen,
		LittleEndian:    conf.Server.Gate.LittleEndian,
		Processor:       protocol.Processor,
		AgentChanRPC:    game.ChanRPC,
	}
}
