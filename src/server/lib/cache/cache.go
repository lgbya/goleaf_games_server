package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var mCache *cache.Cache


type cacheMod struct {

}

//初始化有关全局的缓存变量
func init() {
	mCache = cache.New(5*time.Minute, 10*time.Minute)
	//log.Release("===========启动缓存成功=========== \n")
}

func New() *cacheMod {
	return new(cacheMod)
}

func (c *cacheMod) Set(key string, value interface{}, time time.Duration){
	mCache.Set(key, value, time)
}

func (c *cacheMod) SetNoExpiration(key string, value interface{}) {
	mCache.Set(key, value, cache.NoExpiration)
}


func (c *cacheMod) Get(key string) (interface{}, bool){
	return mCache.Get(key)
}

func (c *cacheMod) Delete(key string)  {
	mCache.Delete(key)
}
