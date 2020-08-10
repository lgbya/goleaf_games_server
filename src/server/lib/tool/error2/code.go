package error2

const (
	//默认错误码，前端只需要显示错误信息
	Default 		= "0"

	//根据前端需要自定义错误码，前端根据错误码进行不同的操作
	ErrSystem       = "100"
	ErrCfg          = "101"
	ErrProtocol     = "102"
)
