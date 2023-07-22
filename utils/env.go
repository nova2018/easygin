package utils

import "github.com/nova2018/easygin/config"

func IsProd() bool {
	env := config.GetString("app_env")
	return env == "prod"
}
