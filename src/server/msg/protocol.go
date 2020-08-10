package msg

type S2C_Error struct {
	Code  string	`json:"code"`
	Message string	`json:"message"`
}

type C2S_Register struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
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
	GameId 	int `json:"game_id"`
}

type S2C_MatchPlayer struct {
	GameId  int `json:"game_id"`
}

type C2S_CancelMatch struct {
}

type S2C_CancelMatch struct {
}

type S2C_StartGame struct {
	RoomId   int                `json:"room_id"`
	UserList map[int]M_UserInfo `json:"user_list"`
}

type C2S_MoraPlaying struct {
	Ply int	`json:"ply"`
}


type S2C_MoraPlaying struct {
	Uid int `json:"uid"`
	Ply int	`json:"ply"`
}

type S2C_MoreReslut struct {
	WinUid int `json:"win_uid"`
}

type M_UserInfo struct {
	Uid 	int 	`json:"uid"`
	Name 	string	`json:"name"`
}