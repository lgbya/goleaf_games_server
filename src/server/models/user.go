package models

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"server/conf"
	"server/lib/cache"
	"server/lib/db"
	"server/lib/tool"
	"time"
)

const (
	GameFree = 0
	GameMath = 1
	GamePlay = 2
)

type User struct {
	ID        int    //id
	Uid       int    //角色id
	Gold      int    //金额
	Name      string //角色名
	Password  string //密码
	Token     string //token
	ExpiresAt int64  //token过期信息
	CreatedAt int64  //创建时间
	Agent 	  *gate.Agent `sql:"-"`
	Game      `sql:"-"`
}

type Agent struct {
	ID 		int //id
}

type Game struct {
	Status			int		//游戏状态
	InRoomId		int			//正在玩的房间id
	GameId 			int 	//所在的游戏id
}

func (u *User) FindLoginName(loginName string) (*User, bool){

	query := db.New().Where("name = ?", loginName).First(u)

	if query.RecordNotFound(){
		return u, false
	}

	if query.Error != nil {
		log.Debug("根据用户名查询出错%v", query.Error)
		return u, false
	}

	return u, true
}


func (u *User) Create(loginName string, loginPassword string) (*User, error) {

	//获取最大的角色id
	newId := new(MaxUid).GetNewUid()

	u.Uid = newId
	u.Gold = 0
	u.Name = loginName
	u.Password = u.getSignPassword(loginPassword)
	u.GenerateToken() //生成token和过期时间
	u.CreatedAt = time.Now().Unix()

	if err := db.New().Create(u).Error; err != nil {
		return u, err
	}

	return u, nil
}

//刷新缓存
//func (u *UserList) FlushCache() {
//	if u.Token != "" && u.Uid > 0 {
//		u.Uid3User(u)
//	}
//}

//删除缓存并修改数据库
func (u *User) DeleteCache(uid int)  {
	user := u.uid4User(uid)
	db.New().Save(user)
}

//删除缓存
func (u *User) uid4User(uid int)  *User {
	if user, ok := u.Uid2User(uid); ok{
		cache.New().Delete(u.ckUid2User(u.Uid))
		return user
	}
	return u
}

//根据uid修改用户信息
func (u *User) Uid3User(user *User)  {
	cache.New().SetNoExpiration(u.ckUid2User(user.Uid), user)
}

//根据uid修改用户角色
func (u *User) Uid2User(uid int)  (*User, bool) {

	if data , found :=cache.New().Get(u.ckUid2User(uid)); found{
		return data.(*User), found
	}
	return nil, false
}

func (u *User) ckUid2User(uid int) string {
	return "key:uid|data:user=" + string(uid)
}

func (u *User)  GenerateToken() (string, int64){
	u.Token = tool.Md5(string(tool.RandNum(999999)) + string(u.ID) + conf.Server.Md5Key)
	u.ExpiresAt = time.Now().AddDate(0,0, 3).Unix()
	return u.Token, u.ExpiresAt
}

func (u *User) AuthLoginPassword(loginPassword string) bool {
	return u.Name == u.getSignPassword(loginPassword)
}

func(u *User)  getSignPassword(password string) string{
	return tool.Md5(password + conf.Server.Md5Key)
}