package app

import (
	"github.com/nova2018/goconfig"
	"github.com/nova2018/goutils"
	"github.com/spf13/viper"
)

type Application interface {
	Run()

	GetCliArgs() *viper.Viper

	GetConfigPath() string

	GetConfig() *goconfig.Config

	RunWithPanic(handle goutils.RecoveryHandle)
}

var (
	_app Application
)

func SetApplication(app Application) {
	_app = app
}

func GetApplication() Application {
	return _app
}
