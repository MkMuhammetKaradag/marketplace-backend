// internal/order-service/config/config.go

package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port        string `mapstructure:"port"`
	GrpcPort    string `mapstructure:"grpcPort"`
	Host        string `mapstructure:"host"`
	Description string `mapstructure:"description"`
}
type DatabaseConfig struct {
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
	Host     string `mapstructure:"host"`
}
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

func Read() Config {
	v := viper.New()

	configDir := getCurrentConfigDir()
	fmt.Println(configDir)

	v.AddConfigPath(configDir)
	v.SetConfigType("yaml")

	files := []string{"server.yaml", "database.yaml"}
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
