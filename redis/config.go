package redis

type redisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
	Prefix   string `mapstructure:"prefix"`
}
