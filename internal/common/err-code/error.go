package errCode

import (
	"github.com/name5566/leaf/gate"
	"server/internal/protocol"
)

func Msg(agent gate.Agent, message string) {
	FatalMsg(agent, Default, message)
}

func FatalMsg(agent gate.Agent, code string, message string) {
	agent.WriteMsg(&protocol.S2C_Error{
		Code:    code,
		Message: message,
	})
}
