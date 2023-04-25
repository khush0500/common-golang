package ravenf

import (
	"fmt"

	appconfig "github.com/kukkar/{{APP_NAME}}/conf"
	"github.com/kukkar/common-golang/pkg/components/raven"
)

func GetFarm(key string) (*raven.Farm, error) {
	if f, err := raven.GetFarm(key); err == nil {
		fmt.Println("farm is already initiated, reusing it")
		return f, nil
	}
	logger := new(raven.DefaultLogger)
	return GetCustomFarm(key, logger)
}

func GetCustomFarm(key string, logger raven.Logger) (*raven.Farm, error) {
	ravenConfig, err := getConfig(logger)
	if err != nil {
		return nil, err
	}
	f, err := raven.InitFarm(key, ravenConfig)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func getConfig(logger raven.Logger) (raven.RavenConfig, error) {
	var rconfig raven.RavenConfig
	config, _ := appconfig.GetAppConfig()

	if config.Raven.Use == string(raven.REDIS_FARM) {
		c := raven.RavenConfig{
			Ftype: raven.REDIS_FARM,
			Config: raven.RedisConf{
				Addr:     config.Raven.Redis.Addr,
				PoolSize: config.Raven.Redis.PoolSize,
			},
			Logger: logger,
		}
		return c, nil
	} else if config.Raven.Use == string(raven.REDIS_CLUSTER_FARM) {
		c := raven.RavenConfig{
			Ftype: raven.REDIS_CLUSTER_FARM,
			Config: raven.RedisConf{
				Addr:     config.Raven.RedisCluster.Addrs,
				PoolSize: config.Raven.RedisCluster.PoolSize,
			},
			Logger: logger,
		}
		return c, nil

	}
	return rconfig, fmt.Errorf("Not a Valid Raven Adapter supplied")
}
