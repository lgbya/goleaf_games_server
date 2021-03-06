package protocol

type S2C_MatchPlayer struct {
	GameId	int `json:"gameId"`
}

type S2C_MoraPlay struct {
	Uid	int `json:"uid"`
	Ply	int `json:"ply"`
}

type C2S_TictactoePlay struct {
	Number	int `json:"number"`
}

type C2S_Heart struct {
}

type C2S_Register struct {
	Name	string `json:"name"`
	Password	string `json:"Password"`
	ConfirmPassword	string `json:"confirmPassword"`
}

type M_UserInfo struct {
	Name	string `json:"name"`
	Uid	int `json:"uid"`
}

type S2C_ResetLogin struct {
	Uid	int `json:"uid"`
	Name	string `json:"name"`
	Gold	int `json:"gold"`
	Token	string `json:"token"`
	ExpiresTime	int64 `json:"expiresTime"`
}

type C2S_MoraPlay struct {
	Ply	int `json:"ply"`
}

type C2S_ResetLogin struct {
	Token	string `json:"token"`
}

type S2C_ContinueGame struct {
	RoomId	int `json:"roomId"`
	GameId	int `json:"gameId"`
	UserList	map[int]M_UserInfo `json:"userList"`
	Continue	map[string]interface{} `json:"continue"`
}

type S2C_EndGame struct {
	WinUid	int `json:"winUid"`
	End	map[string]interface{} `json:"end"`
}

type S2C_Register struct {
	Uid	int `json:"uid"`
	Name	string `json:"name"`
	Gold	int `json:"name"`
	Token	string `json:"token"`
	ExpiresTime	int64 `json:"expiresTime"`
}

type S2C_Login struct {
	ExpiresTime	int64 `json:"expiresTime"`
	Uid	int `json:"uid"`
	Name	string `json:"name"`
	Gold	int `json:"gold"`
	Token	string `json:"token"`
}

type C2S_Login struct {
	Password	string `json:"password"`
	Name	string `json:"name"`
}

type C2S_MatchPlayer struct {
	GameId	int `json:"gameId"`
}

type C2S_CancelMatch struct {
}

type S2C_CancelMatch struct {
}

type S2C_StartGame struct {
	RoomId	int `json:"roomId"`
	GameId	int `json:"gameId"`
	UserList	map[int]M_UserInfo `json:"userList"`
	Start	map[string]interface{} `json:"start"`
}

type S2C_TictactoePlay struct {
	CurrentUid	int `json:"currentUid"`
	Uid	int `json:"uid"`
	Number	int `json:"number"`
}

type S2C_Heart struct {
	Time	int64 `json:"time"`
}

type S2C_Error struct {
	Code	string `json:"code"`
	Message	string `json:"message"`
}

