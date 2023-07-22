package redis

import (
	"github.com/nova2018/easygin/config"
)

func Init() {
	v := config.Sub("redis")
	if v == nil {
		return
	}
	var listCfg map[string]redisConfig
	err := v.Unmarshal(&listCfg)
	if err != nil {
		panic(err)
	}
	for k, cfg := range listCfg {
		initEngine(k, cfg)
	}
}
