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
			GlobalLimit: 1000.0 / 60, GlobalBurst: 1000,
			UserLimit: 1000.0 / 60, UserBurst: 1000,
		},
	}
}

type RoutePolicy struct {
	Permissions int64
}

const (
	PermissionNone                  int64 = 0
	PermissionViewProduct           int64 = 1 << 0
	PermissionAdministrator         int64 = 1 << 62
	PermissionApproveOrRejectSeller int64 = 1 << 24
	PermissionManageRoles           int64 = 1 << 32
	PermissionManageOwnStore        int64 = 1 << 10
)

// GetProtectedRoutes returns the set of routes that require authentication
func GetProtectedRoutes() map[string]RoutePolicy {
	return map[string]RoutePolicy{
		"/users/profile": {
			Permissions: PermissionNone,
		},
		"/users/add-user-role/:user_id": {
			Permissions: PermissionManageRoles | PermissionAdministrator,
		},
		"/users/create-role": {
			Permissions: PermissionManageRoles | PermissionAdministrator,
		},
		"/users/change-password": {
			Permissions: PermissionNone,
		},
		"/users/upload-avatar": {
			Permissions: PermissionNone,
		},
		"/sellers/me": {
			Permissions: PermissionNone,
		},
		"/sellers/onboard": {
			Permissions: PermissionNone,
		},
		"/sellers/approve/:seller_id": {
			Permissions: PermissionApproveOrRejectSeller | PermissionAdministrator,
		},
		"/sellers/reject/:seller_id": {
			Permissions: PermissionApproveOrRejectSeller | PermissionAdministrator,
		},
		"/sellers/upload-store-logo/:seller_id": {
			Permissions: PermissionManageOwnStore,
		},
		"/sellers/upload-store-banner/:seller_id": {
			Permissions: PermissionManageOwnStore,
		},
		"/products/create": {
			Permissions: PermissionManageOwnStore,
		},
		"/products/upload/:product_id": {
			Permissions: PermissionManageOwnStore,
		},
		"/products/category": {
			Permissions: PermissionAdministrator,
		},
		"/products/recommended": {
			Permissions: PermissionNone,
		},
		"/products/product/:product_id": {
			Permissions: PermissionViewProduct,
		},
		"/products/toggle-favorite/:product_id": {
			Permissions: PermissionNone,
		},
		"/products/favorites": {
			Permissions: PermissionNone,
		},
		"/products/update/:product_id": {
			Permissions: PermissionManageOwnStore,
		},
		"/products/delete/:product_id": {
			Permissions: PermissionManageOwnStore | PermissionAdministrator,
		},
		"/baskets/add-item": {
			Permissions: PermissionNone,
		},
		"/baskets/remove-item/:product_id": {
			Permissions: PermissionNone,
		},
		"/baskets/decrement-item/:product_id": {
			Permissions: PermissionNone,
		},
		"/baskets/increment-item/:product_id": {
			Permissions: PermissionNone,
		},
		"/baskets/clear-basket": {
			Permissions: PermissionNone,
		},
		"/baskets/basket": {
			Permissions: PermissionNone,
		},
		"/baskets/count": {
			Permissions: PermissionNone,
		},
	}
}
func Read() Config {
	v := viper.New()

	configDir := getCurrentConfigDir()
	fmt.Println(configDir)

	v.AddConfigPath(configDir)
	v.SetConfigType("yaml")

	files := []string{"server.yaml", "redisCache.yaml"}
	for _, f := range files {
		v.SetConfigFile(filepath.Join(configDir, f))
		if err := v.MergeInConfig(); err == nil {
			fmt.Printf("Config loaded: %s\n", f)
		}
	}

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
