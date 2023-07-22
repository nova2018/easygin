package config

import "github.com/nova2018/goconfig"

var (
	_config *goconfig.Config
)

func setConfig(c *goconfig.Config) {
	_config = c
}

func GetConfig() *goconfig.Config {
	return _config
}
