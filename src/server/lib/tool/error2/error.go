package error2

import (
	"github.com/name5566/leaf/gate"
	"server/msg"
)

func Msg(agent gate.Agent, message string){
	FatalMsg(agent, Default, message)
}

func FatalMsg(agent gate.Agent, code string,  message string){
	agent.WriteMsg(&msg.S2C_Error{
		Code:     code,
		Message:  message,
	})
}


