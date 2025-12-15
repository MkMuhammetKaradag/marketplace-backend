package config

import (
	"time"
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
