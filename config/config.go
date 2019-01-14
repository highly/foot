package config

import "C"
import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

var vipe = New()

type config struct {
	*viper.Viper
}

func New() *config {
	v := &config{viper.New()}
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return v
}

func EnvPrefix(prefix string) {
	vipe.SetEnvPrefix(prefix)
}

func ConfigPath(configPath string) {
	vipe.AddConfigPath(configPath)
}

// [JSON、TOML、YAML、HCL] supported
// config name only with no extension
func Load(configName string) error {
	vipe.SetConfigName(configName)
	return vipe.ReadInConfig()
}

// config file (full path with name and extension)
func LoadFile(configFile string) error {
	vipe.SetConfigFile(configFile)
	return vipe.ReadInConfig()
}

func Watching() {
	vipe.WatchConfig()
	action := func(e fsnotify.Event) {
		log.Println(e.Name, " config has been updated")
	}
	vipe.OnConfigChange(action)
}

func Has(key string) bool {
	return vipe.IsSet(key)
}

func Set(key string, val interface{}) {
	vipe.Set(key, val)
}

func SetDefault(key string, val interface{}) {
	vipe.SetDefault(key, val)
}

func String(key string) string {
	return vipe.GetString(key)
}

func DefaultString(key string, val string) string {
	if Has(key) {
		return String(key)
	}
	return val
}

func Bool(key string) bool {
	return vipe.GetBool(key)
}

func DefaultBool(key string, val bool) bool {
	if Has(key) {
		return Bool(key)
	}
	return val
}

func Int(key string) int {
	return vipe.GetInt(key)
}

func DefaultInt(key string, val int) int {
	if Has(key) {
		return Int(key)
	}
	return val
}

func Int32(key string) int32 {
	return vipe.GetInt32(key)
}

func DefaultInt32(key string, val int32) int32 {
	if Has(key) {
		return Int32(key)
	}
	return val
}

func Int64(key string) int64 {
	return vipe.GetInt64(key)
}

func DefaultInt64(key string, val int64) int64 {
	if Has(key) {
		return Int64(key)
	}
	return val
}

func Float64(key string) float64 {
	return vipe.GetFloat64(key)
}

func DefaultFloat64(key string, val float64) float64 {
	if Has(key) {
		return Float64(key)
	}
	return val
}

func Time(key string) time.Time {
	return vipe.GetTime(key)
}

func DefaultTime(key string, val time.Time) time.Time {
	if Has(key) {
		return Time(key)
	}
	return val
}

func Duration(key string) time.Duration {
	return vipe.GetDuration(key)
}

func DefaultDuration(key string, val time.Duration) time.Duration {
	if Has(key) {
		return Duration(key)
	}
	return val
}

func StringSlice(key string) []string {
	return vipe.GetStringSlice(key)
}

func DefaultStringSlice(key string, val []string) []string {
	if Has(key) {
		return StringSlice(key)
	}
	return val
}

func StringMap(key string) map[string]interface{} {
	return vipe.GetStringMap(key)
}

func DefaultStringMap(key string, val map[string]interface{}) map[string]interface{} {
	if Has(key) {
		return StringMap(key)
	}
	return val
}

func StringMapString(key string) map[string]string {
	return vipe.GetStringMapString(key)
}

func DefaultStringMapString(key string, val map[string]string) map[string]string {
	if Has(key) {
		return StringMapString(key)
	}
	return val
}

func StringMapStringSlice(key string) map[string][]string {
	return vipe.GetStringMapStringSlice(key)
}

func DefaultStringMapStringSlice(key string, val map[string][]string) map[string][]string {
	if Has(key) {
		return StringMapStringSlice(key)
	}
	return val
}

func SizeInBytes(key string) uint {
	return vipe.GetSizeInBytes(key)
}

func DefaultSizeInBytes(key string, val uint) uint {
	if Has(key) {
		return SizeInBytes(key)
	}
	return val
}

func All() map[string]interface{} {
	return vipe.AllSettings()
}
