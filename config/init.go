package config

import (
	"github.com/nova2018/goconfig"
	"github.com/spf13/viper"
)

func Init(configFile, defaultConfigFile string) {
	cfg := viper.New()
	// 发现通过cmd运行程序时，configFile == ''
	if configFile == "" {
		configFile = defaultConfigFile
	}
	cfg.SetConfigFile(configFile)
	err := cfg.ReadInConfig()
	if err != nil {
		panic(err)
	}

	c := goconfig.New()
	c.AddViper(cfg)

	c.StartWait()

	setConfig(c)
}
