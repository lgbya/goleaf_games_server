package protocol

import (
	"github.com/name5566/leaf/network/json"
)

var Processor = json.NewProcessor()

func init() {
	Processor.Register(&C2S_Heart{})
	Processor.Register(&S2C_Heart{})

	Processor.Register(&S2C_Error{})

	Processor.Register(&C2S_Register{})
	Processor.Register(&S2C_Register{})

	Processor.Register(&C2S_Login{})
	Processor.Register(&S2C_Login{})

	Processor.Register(&C2S_ResetLogin{})
	Processor.Register(&S2C_ResetLogin{})

	Processor.Register(&C2S_MatchPlayer{})
	Processor.Register(&S2C_MatchPlayer{})

	Processor.Register(&C2S_CancelMatch{})
	Processor.Register(&S2C_CancelMatch{})

	Processor.Register(&S2C_StartGame{})

	Processor.Register(&S2C_EndGame{})
	Processor.Register(&S2C_ContinueGame{})

	Processor.Register(&C2S_MoraPlay{})
	Processor.Register(&S2C_MoraPlay{})
	Processor.Register(&C2S_TictactoePlay{})
	Processor.Register(&S2C_TictactoePlay{})
}
