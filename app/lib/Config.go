package config

import (
	"github.com/spf13/viper"
	"github.com/adKoch/metric_market/core/kafka_sample/app/lib/log"
)



func LoadConfig(relativeConfigPath string, configFile string, configExtension string) {
	err := LoadViperConfig(relativeConfigPath, configFile, configExtension)
	if err !=nil {
		panic("PANIC: Cannot load config:" + err.Error())
	}
}

func LoadViperConfig(path string, file string, extension string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(file)
	viper.SetConfigType(extension)
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	return
}

func GetConfig(name string) string{
	viper.AutomaticEnv()
	val := viper.GetString(name)
	if val == "" {
		log.Error("Could not find config! Config: " + name)
	}
	return val
}