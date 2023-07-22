package db

import "github.com/nova2018/easygin/config"

func Init() {
	v := config.Sub("database")
	if v == nil {
		return
	}
	var listCfg map[string]databaseConfig
	err := v.Unmarshal(&listCfg)
	if err != nil {
		panic(err)
	}
	for k, cfg := range listCfg {
		initEngine(k, cfg)
	}
}
