package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

var C *config

type config struct {
	*viper.Viper
}

func New() *config {
	C = &config{viper.New()}
	C.AutomaticEnv()
	C.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return C
}

func (c *config) EnvPrefix(prefix string) *config {
	c.SetEnvPrefix(prefix)
	return c
}

func (c *config) ConfigPath(configPath string) *config {
	c.AddConfigPath(configPath)
	return c
}

func (c *config) ConfigName(configName string) *config {
	c.SetConfigName(configName)
	return c
}

// [JSON、TOML、YAML、HCL] supported
// configName only with no extension
func (c *config) Load(configName string) *config {
	if err := c.ConfigName(configName).ReadInConfig(); err != nil {
		panic(fmt.Errorf("Reading config file failed, %s \n", err))
	}
	return c
}

func (c *config) Watching() error {
	c.WatchConfig()
	action := func(e fsnotify.Event) {
		log.Println(e.Name, " config has been updated")
	}
	c.OnConfigChange(action)
	return nil
}

func Set(key string, value interface{}) {
	C.Set(key, value)
}

func Get(key string) interface{} {
	return C.Get(key)
}

func GetString(key string) string {
	return C.GetString(key)
}

func GetBool(key string) bool {
	return C.GetBool(key)
}

func GetInt(key string) int {
	return C.GetInt(key)
}

func GetInt32(key string) int32 {
	return C.GetInt32(key)
}

func GetInt64(key string) int64 {
	return C.GetInt64(key)
}

func GetFloat64(key string) float64 {
	return C.GetFloat64(key)
}

func GetTime(key string) time.Time {
	return C.GetTime(key)
}

func GetDuration(key string) time.Duration {
	return C.GetDuration(key)
}

func GetStringSlice(key string) []string {
	return C.GetStringSlice(key)
}

func GetStringMap(key string) map[string]interface{} {
	return C.GetStringMap(key)
}

func GetStringMapString(key string) map[string]string {
	return C.GetStringMapString(key)
}

func GetStringMapStringSlice(key string) map[string][]string {
	return C.GetStringMapStringSlice(key)
}

func GetSizeInBytes(key string) uint {
	return C.GetSizeInBytes(key)
}
