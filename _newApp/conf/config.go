package config

import (
	"errors"
	"fmt"

	"github.com/kukkar/common-golang/pkg/config"
	"github.com/kukkar/common-golang/pkg/factory/cache"
	"github.com/kukkar/common-golang/pkg/factory/sql"
)

type AppConfig struct {
	MySql              *sql.MysqlConfig   `json:"Mysql"`
	Cache              *cache.CacheConfig `json:"Cache"`
	UseClientValidator bool               `json:"UseClientValidator"`
	SentryDSN          string             `json:"SentryDSN"`
}

func GetAppConfig() (*AppConfig, error) {
	c := config.GlobalAppConfig.ApplicationConfig
	appConfig, ok := c.(*AppConfig)
	if !ok {
		msg := fmt.Sprintf("Example APP Config Not correct %+v", c)
		return nil, errors.New(msg)
	}
	return appConfig, nil
}

func GetGlobalConfig() (*config.AppConfig, error) {
	return config.GlobalAppConfig, nil
}

func EnvUpdateMap() map[string]string {
	m := make(map[string]string)

	m["Mysql.User"] = "{{APP_NAME}}_MYSQL_USER"
	m["Mysql.Password"] = "{{APP_NAME}}_MYSQL_PASSWORD"
	m["Mysql.DbName"] = "{{APP_NAME}}_MYSQL_DBNAME"
	m["Mysql.MaxOpenConnections"] = "{{APP_NAME}}_MYSQL_MAXOPENCONNECTIONS"
	m["Mysql.MaxIdleConnections"] = "{{APP_NAME}}_MYSQL_MAXIDLECONNECTIONS"
	m["Mysql.DefaultTimeZone"] = "{{APP_NAME}}_MYSQL_DEFAULTTIMEZONE"
	m["Mysql.Host"] = "{{APP_NAME}}_MYSQL_HOST"
	m["Mysql.Port"] = "{{APP_NAME}}_MYSQL_PORT"
	m["Cache.Use"] = "{{APP_NAME}}_CACHE_USE"
	m["Cache.Redis.Addr"] = "{{APP_NAME}}_CACHE_REDIS_ADDRESS"
	m["Cache.Redis.PoolSize"] = "{{APP_NAME}}_CACHE_REDIS_POOLSIZE"

	m["SentryDSN"] = "{{APP_NAME}}_SENTRY_DSN"
	return m
}
