package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"os"
)

var Cfg *Config

type Config struct {
	Authorization []string `mapstructure:"Authorization" validate:"required"`
}

// LoadConfig from config file .
func LoadConfig(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cumulx" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cumulx")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("read config failed:%v:", err))
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		panic(fmt.Sprintf("unmarshal config to object failed:%v:", err))
	}

}