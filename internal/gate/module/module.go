package module

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/module"
	"server/internal/common/conf"
	"server/internal/game"
	"server/internal/gate/protocol"
	"time"
)

type Module struct {
	*gate.Gate
}

var _ module.Module = (*Module)(nil)
func (m *Module) OnInit() {
	config := conf.Get()

	m.Gate = &gate.Gate{
		MaxConnNum:      config.Server.MaxConnNum,
		WSAddr:          config.Server.WSAddr,
		CertFile:        config.Server.CertFile,
		KeyFile:         config.Server.KeyFile,
		TCPAddr:         config.Server.TCPAddr,

		HTTPTimeout:     config.Gate.HTTPTimeout * time.Second,
		PendingWriteNum: config.Gate.PendingWriteNum,
		MaxMsgLen:       config.Gate.MaxMsgLen,
		LenMsgLen:       config.Gate.LenMsgLen,
		LittleEndian:    config.Gate.LittleEndian,

		Processor:       protocol.Processor,

		AgentChanRPC:    game.ChanRPC,
	}
}
