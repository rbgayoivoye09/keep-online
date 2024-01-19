package config

import (
	"fmt"
	"path/filepath"

	. "github.com/rbgayoivoye09/keep-online/src/utils/log"

	"github.com/spf13/viper"
)

// Config 结构体用于存储配置信息

type User struct {
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
}
type Mail struct {
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
	Server   string `mapstructure:"server"`
}

type Web struct {
	RedirURL string `mapstructure:"redirurl"`
	LoginURL string `mapstructure:"login_url"`
}

type SSH struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	FilePath string `mapstructure:"filePath"`
}

type Config struct {
	User `mapstructure:"user"`

	Mail `mapstructure:"mail"`

	Web `mapstructure:"web"`

	SSH `mapstructure:"ssh"`
}

// NewConfig 用于创建一个新的 Config 实例
func NewConfig(configFilePath string) (*Config, error) {
	viper.SetConfigFile(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}

var config *Config

func GetConfig() *Config {
	return config
}

func init() {

	projectRoot, err := filepath.Abs(".") // 根据你的项目结构调整相对路径
	if err != nil {
		Logger.Sugar().Fatalf("Error getting project root path: %v", err)
	}
	configFilePath := filepath.Join(projectRoot, "config", "user.yml")
	config, err = NewConfig(configFilePath)
	if err != nil {
		Logger.Sugar().Fatalf("Error reading config file: %v", err)
	}

}
