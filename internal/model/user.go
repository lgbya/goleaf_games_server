package model

import (
	"server/internal/common/cache"
	"server/internal/common/conf"
	"server/internal/common/db"
	"server/internal/common/helper"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"

	"time"
)

type User struct {
	ID        int         //id
	Uid       int         //角色id
	Gold      int         //金额
	Name      string      //角色名
	Password  string      //密码
	Token     string      //token
	ExpiresAt int64       //token过期信息
	CreatedAt int64       //创建时间
	IsRobot   bool        `sql:"-"`
	Agent     *gate.Agent `sql:"-"`
	Game      Game        `sql:"-"`
}

type Agent struct {
	ID        int   //id
	HeartTime int64 //收到心跳包的时间
}

type Game struct {
	Id       int //所在的游戏id
	Status   int //游戏状态
	InRoomId int //正在玩的房间id
}

func (u User) FindLoginName(loginName string) (*User, bool) {

	user := &User{}
	query := db.New().Where("name = ?", loginName).First(user)

	if query.RecordNotFound() {
		return nil, false
	}

	if query.Error != nil {
		log.Debug("根据用户名查询出错%v", query.Error)
		return nil, false
	}

	return user, true
}

func (u *User) Create(loginName string, loginPassword string) error {

	//获取最大的角色id
	newId := new(MaxUid).GetNewUid()

	u.Uid = newId
	u.Gold = 0
	u.Name = loginName
	u.Password = u.getSignPassword(loginPassword)
	u.GenerateToken() //生成token和过期时间
	u.CreatedAt = time.Now().Unix()

	if err := db.New().Create(u).Error; err != nil {
		return  err
	}

	return nil
}

//删除缓存并修改数据库
func (u User) DeleteCache(uid int) {
	user := u.uid4User(uid)
	if user.IsRobot == false {
		db.New().Save(user)
	}
}

//用户uid对应的用户信息
func (u User) ckUid2User(uid int) string {
	return "key:uid|data:user=" + string(uid)
}

func (u User) Uid2User(uid int) (*User, bool) {

	if data, found := cache.New().Get(u.ckUid2User(uid)); found {
		return data.(*User), found
	}
	return nil, false
}

func (u User) Uid3User(user *User) {
	cache.New().SetNoExpiration(u.ckUid2User(user.Uid), user)
}

func (u User) uid4User(uid int) *User {
	if user, ok := u.Uid2User(uid); ok {
		cache.New().Delete(u.ckUid2User(u.Uid))
		return user
	}
	return nil
}

//登录后存储在公共map中
func (u User) ckCommon2LoginUid() string {
	return "common|data:uid_list"
}

func (u User) Common2LoginUid() map[int]*gate.Agent {
	if data, found := cache.New().Get(u.ckCommon2LoginUid()); found {
		return data.(map[int]*gate.Agent)
	}
	return make(map[int]*gate.Agent)
}

func (u User) Common3LoginUid(uid int, agent *gate.Agent) {
	uid2Agent := u.Common2LoginUid()
	uid2Agent[uid] = agent
	cache.New().SetNoExpiration(u.ckCommon2LoginUid(), uid2Agent)
}

func (u User) Common4LoginUid(uid int) {
	uid2Agent := u.Common2LoginUid()
	delete(uid2Agent, uid)
	cache.New().SetNoExpiration(u.ckCommon2LoginUid(), uid2Agent)
}

//临时的存储token对应的用户关系
func (u User) ckTempToken2User(token string) string {
	return "key:uid|data:reset_user=" + token
}

func (u User) TempToken2User(token string) (*User, bool) {

	if data, found := cache.New().Get(u.ckTempToken2User(token)); found {
		return data.(*User), found
	}
	return nil, false
}

func (u User) TempToken3User(user *User, time time.Duration) {
	cache.New().Set(u.ckTempToken2User(user.Token), user, time)
}

func (u User) TempToken4User(token string) *User {
	if user, ok := u.TempToken2User(token); ok {
		cache.New().Delete(u.ckUid2User(u.Uid))
		return user
	}
	return nil
}

func (u User) CheckRepeatLogin(uid int) bool {
	uid2Agent := u.Common2LoginUid()
	_, ok := uid2Agent[uid]
	return ok
}

func (u *User) GenerateToken() (string, int64) {
	u.Token = helper.Md5(string(helper.RandNum(999999)) + string(u.ID) + conf.Get().Server.Md5Key)
	u.ExpiresAt = time.Now().AddDate(0, 0, 3).Unix()
	return u.Token, u.ExpiresAt
}

func (u User) AuthLoginPassword(loginPassword string) bool {
	return u.Password == u.getSignPassword(loginPassword)
}

func (u User) getSignPassword(password string) string {
	return helper.Md5(password + conf.Get().Server.Md5Key)
}

//写入登录数据
func (u User) SetLoginInfo(user *User, agent gate.Agent) *User {
	user.GenerateToken()
	user.Agent = &agent
	agent.SetUserData(&Agent{ID: user.Uid, HeartTime: time.Now().Unix()})
	user.Uid3User(user)
	user.Common3LoginUid(user.Uid, &agent)
	return user
}
