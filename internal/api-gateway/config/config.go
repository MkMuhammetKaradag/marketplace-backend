package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Configuration Constants
const (
	// InternalGatewayHeader is the header key used for internal service communication
	InternalGatewayHeader = "X-API-Key"
	// InternalGatewaySecret should ideally be loaded from environment variables
	InternalGatewaySecret = "GATEWAY_SECRET_KEY"

	SessionCookieName = "Session"
	DefaultTimeout    = 30 * time.Second
	GatewayPort       = ":8080"
)

// RouteConfig defines rate limits for specific routes
type RouteConfig struct {
	GlobalLimit float64 // Requests per second
	GlobalBurst int
	UserLimit   float64 // Requests per second
	UserBurst   int
}
type ServerConfig struct {
	Port        string `mapstructure:"port"`
	GrpcPort    string `mapstructure:"grpcPort"`
	Host        string `mapstructure:"host"`
	Description string `mapstructure:"description"`
}
type RedisCacheConfig struct {
	Addr     string        `mapstructure:"addr"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db"`
	CacheTTL time.Duration `mapstructure:"cache_ttl"`
}

type Config struct {
	RedisCache RedisCacheConfig `mapstructure:"redisCache"`
	Server     ServerConfig     `mapstructure:"server"`
}

// GetDefaultRouteConfigs returns the default rate limit configurations
func GetDefaultRouteConfigs() map[string]RouteConfig {
	return map[string]RouteConfig{
		"/users": {
			GlobalLimit: 50.0 / 60, GlobalBurst: 50,
			UserLimit: 20.0 / 60, UserBurst: 20,
		},
		// "/users/profile": {
		// 	GlobalLimit: 10.0 / 60, GlobalBurst: 10,
		// 	UserLimit: 2.0 / 60, UserBurst: 2,
		// },
		"/sellers": {
			GlobalLimit: 50.0 / 60, GlobalBurst: 50,
			UserLimit: 20.0 / 60, UserBurst: 20,
		},
		"/test/hello": {
			GlobalLimit: 2.0 / 60, GlobalBurst: 2,
			UserLimit: 1.0 / 60, UserBurst: 1,
		},
		"/chat": {
			GlobalLimit: 50.0 / 60, GlobalBurst: 10,
			UserLimit: 5.0 / 60, UserBurst: 3,
		},
		"default": {
			GlobalLimit: 1.0 / 60, GlobalBurst: 1,
			UserLimit: 1.0 / 60, UserBurst: 1,
		},
	}
}

type RoutePolicy struct {
	Roles []string
}

// GetProtectedRoutes returns the set of routes that require authentication
func GetProtectedRoutes() map[string]RoutePolicy {
	return map[string]RoutePolicy{
		"/users/profile": {
			Roles: []string{"buyer", "seller", "admin"},
		},
		"/users/change-user-role/:user_id": {
			Roles: []string{"admin"},
		},
		"/sellers/onboard": {
			Roles: []string{"buyer", "admin"},
		},
		"/sellers/approve/:seller_id": {
			Roles: []string{"admin"},
		},
		"/sellers/reject/:seller_id": {
			Roles: []string{"admin"},
		},
	}
}
func Read() Config {
	v := viper.New()

	// Bu dosyanın kendi klasörünü al (internal/user-service/config)
	configDir := getCurrentConfigDir()
	fmt.Println(configDir)

	v.AddConfigPath(configDir) // artık kesin doğru yer
	v.SetConfigType("yaml")

	// Dosyaları sırayla yükle (varsa)
	files := []string{"server.yaml", "redisCache.yaml"}
	for _, f := range files {
		v.SetConfigFile(filepath.Join(configDir, f))
		if err := v.MergeInConfig(); err == nil {
			fmt.Printf("Config loaded: %s\n", f)
		}
	}

	// ENV override (en son gelir)
	v.SetEnvPrefix("USER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic("Config unmarshal error: " + err.Error())
	}

	return cfg
}

// Bu fonksiyon bu dosyanın bulunduğu klasörü döndürür
func getCurrentConfigDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("Config folder not found")
	}
	return filepath.Dir(file) // ← bu dosyanın olduğu klasör: internal/user-service/config
}
