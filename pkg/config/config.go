package config

import (
	"time"

	"github.com/spf13/viper"
)

type config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	} `yaml:"database"`
	Log struct {
		Level string `yaml:"level"`
		Path  string `yaml:"path"`
	} `yaml:"log"`
	Jwt struct {
		SecretKey      string        `yaml:"secretKey"`
		AccessTimeout  time.Duration `yaml:"accessTimeout"`
		RefreshTimeout time.Duration `yaml:"refreshTimeout"`
	} `yaml:"jwt"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
	Snowflake struct {
		NodeId int64 `yaml:"nodeId"`
	} `yaml:"snowflake"`
}

var instance *config

func InitConfig() error {
	var err error
	instance = &config{}
	err = instance.load("config/config.yaml")
	return err
}

func GetConfig() *config {
	return instance
}

func (c *config) load(path string) error {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}
