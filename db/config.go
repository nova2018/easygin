package db

import (
	"github.com/nova2018/easygin/file"
	"time"
)

type databaseConfig struct {
	Driver        string            `mapstructure:"driver"`
	Connection    string            `mapstructure:"connection"`
	Prefix        string            `mapstructure:"prefix"`
	SlowThreshold time.Duration     `mapstructure:"slow_threshold"`
	ConsoleLog    bool              `mapstructure:"console_log"`
	SqlLog        string            `mapstructure:"sql_log"`
	LogConfig     *file.WriteConfig `mapstructure:"log_config"`
	LogLevel      int               `mapstructure:"log_level"`
	Pool          *struct {
		MaxIdleConns    int           `mapstructure:"max_idle_conns"`
		MaxOpenConns    int           `mapstructure:"max_open_conns"`
		ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	} `mapstructure:"pool"`
}
