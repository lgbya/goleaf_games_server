package msg

type C2S_Heart struct {
}

type S2C_Heart struct {
	Time 	int64 	`json:"time"`
}

type S2C_Error struct {
	Code  string	`json:"code"`
	Message string	`json:"message"`
}

type C2S_Register struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}


type S2C_Register struct {
	Uid 	int 	`json:"uid"`
	Name 	string	`json:"name"`
	Gold	int 	`json:"gold"`
	Token 	string	`json:"token"`
	ExpiresTime int64 `json:"token"`
}

type C2S_Login struct {
	Name     string	`json:"name"`
	Password string	`json:"password"`
}

type S2C_Login struct {
	Uid 	int 	`json:"uid"`
	Name 	string	`json:"name"`
	Gold	int 	`json:"gold"`
	Token 	string	`json:"token"`
	ExpiresTime int64 `json:"token"`
}

type C2S_MatchPlayer struct {
	GameId 	int `json:"gameId"`
}

type S2C_MatchPlayer struct {
	GameId  int `json:"gameId"`
}

type C2S_CancelMatch struct {
}

type S2C_CancelMatch  struct {
}

type S2C_StartGame struct {
	RoomId   int                `json:"roomId"`
	GameId 	 int 				`json:"gameId"`
	UserList map[int]M_UserInfo `json:"userList"`
	Start map[string]interface{} `json:"start"`
}

type S2C_ContinueGame struct {
	RoomId   int                `json:"roomId"`
	GameId 	 int 				`json:"gameId"`
	UserList map[int]M_UserInfo `json:"userList"`
	Continue map[string]interface{} `json:"continue"`
}

type C2S_MoraPlaying struct {
	Ply int	`json:"ply"`
}


type S2C_MoraPlaying struct {
	Uid int `json:"uid"`
	Ply int	`json:"ply"`
}

type S2C_MoreResult struct {
	WinUid int `json:"win_uid"`
	GameInfo map[int]int `json:"gameInfo"`
}

type M_UserInfo struct {
	Uid 	int 	`json:"uid"`
	Name 	string	`json:"name"`
}