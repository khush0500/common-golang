package cache

import (
	"fmt"

	"github.com/kukkar/common-golang/pkg/factory/cache"

	"github.com/go-redis/redis"
	concurrenthashmap "github.com/kukkar/common-golang/pkg/utils/concurrenthashmap"
	appConf "github.com/kukkar/{{APP_NAME}}}/conf"
)

// DefaultKey default pool key for mysql conn
const DefaultKey = "default"

var cacheMap = concurrenthashmap.New()

func GetPool(key string) (*redis.Client, error) {
	if val, ok := cacheMap.Get(key); !ok {
		//we dont have a pool by this key, initiate new pool.
		pool, err := InitPool(key)
		if err != nil {
			return nil, fmt.Errorf("Could not initiate pool for key:%s, Error:%s",
				key, err.Error())
		}
		cacheMap.Put(key, pool)
		return pool, nil
	} else {
		return val.(*redis.Client), nil
	}
}

func InitPool(key string) (*redis.Client, error) {
	conf, err := appConf.GetAppConfig()
	if err != nil {
		return nil, err
	}
	if conf.Cache == nil {
		return nil, fmt.Errorf("cache config can not be empty")
	}
	return cache.InitPool((*conf.Cache))
}
