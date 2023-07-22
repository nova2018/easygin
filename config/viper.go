package config

import (
	"github.com/spf13/viper"
	"time"
)

func GetViper() *viper.Viper {
	return GetConfig().GetConfig()
}

func Get(key string) interface{} { return GetViper().Get(key) }

func Sub(key string) *viper.Viper { return GetViper().Sub(key) }

func GetString(key string) string { return GetViper().GetString(key) }

func GetBool(key string) bool { return GetViper().GetBool(key) }

func GetInt(key string) int { return GetViper().GetInt(key) }

func GetInt32(key string) int32 { return GetViper().GetInt32(key) }

func GetInt64(key string) int64 { return GetViper().GetInt64(key) }

func GetUint(key string) uint { return GetViper().GetUint(key) }

func GetUint16(key string) uint16 { return GetViper().GetUint16(key) }

func GetUint32(key string) uint32 { return GetViper().GetUint32(key) }

func GetUint64(key string) uint64 { return GetViper().GetUint64(key) }

func GetFloat64(key string) float64 { return GetViper().GetFloat64(key) }

func GetTime(key string) time.Time { return GetViper().GetTime(key) }

func GetDuration(key string) time.Duration { return GetViper().GetDuration(key) }

func GetIntSlice(key string) []int { return GetViper().GetIntSlice(key) }

func GetStringSlice(key string) []string { return GetViper().GetStringSlice(key) }

func GetStringMap(key string) map[string]interface{} { return GetViper().GetStringMap(key) }

func GetStringMapString(key string) map[string]string { return GetViper().GetStringMapString(key) }

func GetStringMapStringSlice(key string) map[string][]string {
	return GetViper().GetStringMapStringSlice(key)
}

func GetSizeInBytes(key string) uint { return GetViper().GetSizeInBytes(key) }

func UnmarshalKey(key string, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return GetViper().UnmarshalKey(key, rawVal, opts...)
}

func Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return GetViper().Unmarshal(rawVal, opts...)
}

func UnmarshalExact(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return GetViper().UnmarshalExact(rawVal, opts...)
}

func IsSet(key string) bool { return GetViper().IsSet(key) }

func InConfig(key string) bool { return GetViper().InConfig(key) }
