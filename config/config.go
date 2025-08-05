package config

import (
	"github.com/spf13/viper"
)

func LoadViperConfigs(paths ...string) *viper.Viper {
	conf := viper.New()

	for _, path := range paths {
		pathConf := viper.New()
		pathConf.SetConfigFile(path)
		if err := pathConf.ReadInConfig(); err != nil {
			panic(err)
		}
		conf.MergeConfigMap(pathConf.AllSettings())
	}

	return conf
}

func NewViperFromMap(ms ...map[string]any) *viper.Viper {
	conf := viper.New()

	for _, m := range ms {
		if m == nil {
			continue
		}
		if err := conf.MergeConfigMap(m); err != nil {
			panic(err)
		}
	}

	return conf
}
