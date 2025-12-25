// internal/user-service/config/config.go

package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
	Host     string `mapstructure:"host"`
}

type CloudinaryConfig struct {
	CloudName string `mapstructure:"cloudName"`
	APIKey    string `mapstructure:"apiKey"`
	APISecret string `mapstructure:"apiSecret"`
}
type ServerConfig struct {
	Port        string `mapstructure:"port"`
	GrpcPort    string `mapstructure:"grpcPort"`
	Host        string `mapstructure:"host"`
	Description string `mapstructure:"description"`
}
type MessagingConfig struct {
	Brokers []string `mapstructure:"brokers"`
}
type RedisSessionConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Config struct {
	Database     DatabaseConfig     `mapstructure:"database"`
	Server       ServerConfig       `mapstructure:"server"`
	RedisSession RedisSessionConfig `mapstructure:"redisSession"`
	Messaging    MessagingConfig    `mapstructure:"messaging"`
	Cloudinary   CloudinaryConfig   `mapstructure:"cloudinary"`
}

func Read() Config {
	v := viper.New()

	// Bu dosyanın kendi klasörünü al (internal/user-service/config)
	configDir := getCurrentConfigDir()
	fmt.Println(configDir)

	v.AddConfigPath(configDir) // artık kesin doğru yer
	v.SetConfigType("yaml")

	// Dosyaları sırayla yükle (varsa)
	files := []string{"server.yaml", "database.yaml", "messasing.yaml", "cloudinary.yaml"}
	for _, f := range files {
		v.SetConfigFile(filepath.Join(configDir, f))
		if err := v.MergeInConfig(); err == nil {
			fmt.Printf("Config yüklendi: %s\n", f)
		}
	}

	// ENV override (en son gelir)
	v.SetEnvPrefix("USER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Varsayılanlar
	v.SetDefault("server.port", "8081")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("database.port", "5432")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "password")
	v.SetDefault("database.db", "marketplace")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic("Config unmarshal hatası: " + err.Error())
	}

	return cfg
}

// Bu fonksiyon bu dosyanın bulunduğu klasörü döndürür
func getCurrentConfigDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("Config klasörü tespit edilemedi")
	}
	return filepath.Dir(file) // ← bu dosyanın olduğu klasör: internal/user-service/config
}
