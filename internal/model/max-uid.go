package model

import (
	"server/internal/common/cache"
	"server/internal/common/db"

	"github.com/name5566/leaf/log"
)

const CkCommon2MaxUid = "key:common|data:max_uid"
const InitialUid = 10000000

type MaxUid struct {
	ID  int `gorm:"primary_key"`
	Max int `gorm:"type:int(10) unsigned;not null;default 0"`
}

func init() {
	maxId, err := new(MaxUid).getDbMax()
	if err != nil {
		log.Fatal("初始化MaxUserId严重错误，程序退出 %v", err)
	}
	cache.New().SetNoExpiration(CkCommon2MaxUid, maxId)
	//log.Release("======MaxUid写入缓存成功======")

}

//获取最新的角色id
func (m *MaxUid) GetNewUid() int {

	newId := InitialUid
	c := cache.New()

	if maxId, ok := c.Get(CkCommon2MaxUid); ok {
		//如果缓存存在最大id,直接取缓存的
		newId = maxId.(int)
	} else {
		//缓存不存在取数据库
		newId, _ = m.getDbMax()
	}

	c.SetNoExpiration(CkCommon2MaxUid, newId+1)

	return newId
}

//获取max_uid的最大值
func (m *MaxUid) getDbMax() (int, error) {

	var model = MaxUid{Max: InitialUid}
	var err error

	//数据库也找不到，插入新的数据
	if db.New().Order("max desc").First(&model).RecordNotFound() {
		err = db.New().Create(&model).Error
		if err != nil {
			log.Error("严重错误，角色最大id表插入失败 %v", err)
		}
	} else if db.New().Error != nil {
		//判断是否查询是否出了其他错误
		err = db.New().Error
	}

	return model.Max, err
}

func (m *MaxUid) FlushDb() {
	if maxUid, found := cache.New().Get(CkCommon2MaxUid); found {

		m.Max = maxUid.(int)
		err := db.New().Model(m).Update("max", maxUid).Error

		if err == nil {
			log.Release("刷新max_uid表成功")
		} else {
			log.Release("刷新max_uid表失败 max_uid:%v, err:%v", maxUid, err)
		}
	}
}
