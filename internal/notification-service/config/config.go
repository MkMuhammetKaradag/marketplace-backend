// internal/notification-service/config/config.go

package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type MessagingConfig struct {
	Brokers []string `mapstructure:"brokers"`
}
type ServerConfig struct {
	Port        string `mapstructure:"port"`
	GrpcPort    string `mapstructure:"grpcPort"`
	Host        string `mapstructure:"host"`
	Description string `mapstructure:"description"`
}

type EmailConfig struct {
	ApiKey string `mapstructure:"apiKey"`
}

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Messaging MessagingConfig `mapstructure:"messaging"`
	Email     EmailConfig     `mapstructure:"email"`
}

func Read() Config {
	v := viper.New()

	configDir := getCurrentConfigDir()
	fmt.Println(configDir)

	v.AddConfigPath(configDir)
	v.SetConfigType("yaml")

	files := []string{"server.yaml", "messaging.yaml", "email.yaml"}
	for _, f := range files {
		v.SetConfigFile(filepath.Join(configDir, f))
		if err := v.MergeInConfig(); err == nil {
			fmt.Printf("Config loaded: %s\n", f)
		}
	}

	v.SetEnvPrefix("USER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("server.port", "8083")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "password")
	v.SetDefault("database.db", "marketplace")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic("Config unmarshal error: " + err.Error())
	}

	return cfg
}

func getCurrentConfigDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("Config folder not found")
	}
	return filepath.Dir(file)
}
