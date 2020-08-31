package cache

import (
	"github.com/patrickmn/go-cache"
	"server/internal/common/conf"
	"time"
)

var _cache *cache.Cache

type cacheMod struct {
}

//初始化有关全局的缓存变量
func init() {
	cacheCfg := conf.Server.Cache
	_cache = cache.New(cacheCfg.DefaultExpiration*time.Minute, cacheCfg.CleanupInterval*time.Minute)
	//log.Release("===========启动缓存成功=========== \n")
}

func New() *cacheMod {
	return new(cacheMod)
}

func (c *cacheMod) Set(key string, value interface{}, time time.Duration) {
	_cache.Set(key, value, time)
}

func (c *cacheMod) SetNoExpiration(key string, value interface{}) {
	_cache.Set(key, value, cache.NoExpiration)
}

func (c *cacheMod) Get(key string) (interface{}, bool) {
	return _cache.Get(key)
}

func (c *cacheMod) Delete(key string) {
	_cache.Delete(key)
}
