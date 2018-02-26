package cache

import "github.com/o3labs/openpoint/platform/cache/local"

type CacheService struct {
	CacheServiceManager
}

type CacheServiceManager interface {
	Get(key string) interface{}
	GetScan(key string, v interface{}) error
	Set(key string, value interface{})
}

func Shared() *CacheService {
	//always local for now
	client := local.NewLocalClient()
	return &CacheService{CacheServiceManager: client}
	// if config.Env.Name == config.LocalEnv {
	// 	client := local.NewLocalClient()
	// 	return &CacheService{CacheServiceManager: client}
	// }

	// //default is production one
	// client := redis.NewRedisClient()
	// return &CacheService{CacheServiceManager: client}
}
