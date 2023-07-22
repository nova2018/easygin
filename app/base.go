package app

import (
	"flag"
	"fmt"
	"github.com/nova2018/easygin/config"
	"github.com/nova2018/easygin/db"
	"github.com/nova2018/easygin/logger"
	"github.com/nova2018/easygin/redis"
	"github.com/nova2018/goconfig"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

var (
	_defaultConfigFile string = "./config/dev.toml"
)

func SetDefaultConfigFile(s string) {
	_defaultConfigFile = s
}

func GetDefaultConfigFile() string {
	return _defaultConfigFile
}

type BaseApplication struct {
	viper *viper.Viper
}

func (b *BaseApplication) init() {
	flag.String("config", GetDefaultConfigFile(), "config")
}

func (b *BaseApplication) parse(f *pflag.FlagSet) {
	if f == nil {
		flag.Usage = func() {
			_, _ = fmt.Fprintf(os.Stderr, "Usage of params:\n")
			flag.PrintDefaults()
		}

		pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
		pflag.Parse()
		f = pflag.CommandLine
	}

	v := viper.New()
	err := v.BindPFlags(f)
	if err != nil {
		panic(err)
	}

	b.viper = v
}

func (b *BaseApplication) GetCliArgs() *viper.Viper {
	return b.viper
}

func (b *BaseApplication) GetConfigPath() string {
	return b.GetCliArgs().GetString("config")
}

func (*BaseApplication) GetConfig() *goconfig.Config {
	return config.GetConfig()
}

func (b *BaseApplication) startUp() {
	// init config
	config.Init(b.GetConfigPath(), GetDefaultConfigFile())

	// init alarm
	//utils.SetConfig(config.GetConfig())

	// init logger
	logger.Init()

	// init mysql
	db.Init()

	// init redis
	redis.Init()

}
