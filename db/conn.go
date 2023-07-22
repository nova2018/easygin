package db

import (
	"fmt"
	"gorm.io/gorm"
)

func Default() *gorm.DB {
	return Connection("default")
}

var (
	_mapDb     map[string]*gorm.DB
	_mapPrefix map[string]string
)

func Connection(name string) *gorm.DB {
	if e, ok := _mapDb[name]; ok {
		return e
	}
	return nil
}

func PrefixDefault() string {
	return Prefix("default")
}

func Prefix(name string) string {
	if x, ok := _mapPrefix[name]; ok {
		return x
	}
	return ""
}

func SetPrefix(name, prefix string) {
	_mapPrefix[name] = prefix
}

func setConnection(name string, engine *gorm.DB) {
	fmt.Printf("mysql connect success! name=[%s]\n", name)
	_mapDb[name] = engine
}

func init() {
	_mapDb = make(map[string]*gorm.DB, 1)
	_mapPrefix = make(map[string]string, 1)
}
