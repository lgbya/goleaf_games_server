package cache

import (
	"github.com/patrickmn/go-cache"
	"server/internal/common/conf"
	"time"
)

var _cache *cache.Cache

type Cache struct {
}

//初始化有关全局的缓存变量
func init() {
	config := conf.Get().Cache
	_cache = cache.New(config.DefaultExpiration*time.Minute, config.CleanupInterval*time.Minute)
	//log.Release("===========启动缓存成功=========== \n")
}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) Set(key string, value interface{}, time time.Duration) {
	_cache.Set(key, value, time)
}

func (c *Cache) SetNoExpiration(key string, value interface{}) {
	_cache.Set(key, value, cache.NoExpiration)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return _cache.Get(key)
}

func (c *Cache) Delete(key string) {
	_cache.Delete(key)
}
