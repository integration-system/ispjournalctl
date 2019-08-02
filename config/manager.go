package config

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultConfigPath = "/etc/mdm-dump/config.yml"
	devConfigPath     = "./config/config.yml"
)

var (
	config Configuration
)

func Get() Configuration {
	return config
}

func Load(path string) (Configuration, error) {
	if path != "" {
		return extractConfig(path)
	} else if os.Getenv("APP_MODE") == "dev" {
		return extractConfig(devConfigPath)
	} else {
		return extractConfig(defaultConfigPath)
	}
}

func extractConfig(pathConf string) (Configuration, error) {
	viper.AutomaticEnv()

	configPath := filepath.Dir(pathConf)

	configFileName := filepath.Base(pathConf)
	extension := filepath.Ext(pathConf)
	configFileName = strings.Replace(configFileName, extension, "", 1)

	viper.AddConfigPath(configPath)
	viper.SetConfigName(configFileName)

	return loadConfig()
}

func loadConfig() (Configuration, error) {
	cfg := Configuration{}

	if err := viper.ReadInConfig(); err != nil {
		return cfg, errors.Wrap(err, "read config file")
	} else if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, errors.Wrap(err, "unmarshal config")
	} else if _, err := govalidator.ValidateStruct(cfg); err != nil {
		m := govalidator.ErrorsByField(err)
		return cfg, fmt.Errorf("invalid config: %v", m)
	}

	config = cfg
	return cfg, nil
}
