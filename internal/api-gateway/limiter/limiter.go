package limiter

import (
	"log"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type LimiterEntry struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

type RateLimiter struct {
	limiters sync.Map // map[string]*LimiterEntry
	stopChan chan struct{}
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		stopChan: make(chan struct{}),
	}
}

func (rl *RateLimiter) GetLimiter(key string, r rate.Limit, b int) *rate.Limiter {
	now := time.Now()
	if entry, ok := rl.limiters.Load(key); ok {
		limiterEntry := entry.(*LimiterEntry)
		limiterEntry.lastAccess = now
		return limiterEntry.limiter
	}
	newLimiter := rate.NewLimiter(r, b)
	entry := &LimiterEntry{
		limiter:    newLimiter,
		lastAccess: now,
	}
	actual, loaded := rl.limiters.LoadOrStore(key, entry)
	if loaded {
		return actual.(*LimiterEntry).limiter
	}
	return newLimiter
}

// func (rl *RateLimiter) StartCleanup(interval, maxAge time.Duration) {
// 	ticker := time.NewTicker(interval)
// 	go func() {
// 		for range ticker.C {
// 			count := 0
// 			now := time.Now()
// 			rl.limiters.Range(func(key, value interface{}) bool {
// 				entry := value.(*LimiterEntry)
// 				if now.Sub(entry.lastAccess) > maxAge {
// 					rl.limiters.Delete(key)
// 					count++
// 				}
// 				return true
// 			})
// 			if count > 0 {
// 				log.Printf("🧹 Cleanup: %d limiters removed", count)
// 			}
// 		}
// 	}()
// }

func (rl *RateLimiter) StartCleanup(interval, maxAge time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				rl.cleanup(maxAge)
			case <-rl.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}
func (rl *RateLimiter) cleanup(maxAge time.Duration) {
	count := 0
	now := time.Now()

	rl.limiters.Range(func(key, value interface{}) bool {
		entry := value.(*LimiterEntry)
		if now.Sub(entry.lastAccess) > maxAge {
			rl.limiters.Delete(key)
			count++
		}
		return true
	})

	if count > 0 {
		log.Printf("🧹 Temizlik: %d limiter silindi", count)
	}
}

func (rl *RateLimiter) Stop() {
	close(rl.stopChan)
}

// RouteConfig her yol için rate limit yapılandırmasını tutar
type RouteConfig struct {
	GlobalLimit rate.Limit
	GlobalBurst int
	UserLimit   rate.Limit
	UserBurst   int
}

type ConfigManager struct {
	configs map[string]RouteConfig
	mu      sync.RWMutex
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		configs: make(map[string]RouteConfig),
	}
}

func (cm *ConfigManager) SetConfig(path string, config RouteConfig) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.configs[path] = config
}

func (cm *ConfigManager) GetConfig(path string) RouteConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Tam eşleşme kontrolü
	if config, exists := cm.configs[path]; exists {
		return config
	}

	// Prefix eşleşmesi
	longestMatch := ""
	var matchedConfig RouteConfig
	foundMatch := false

	for route, config := range cm.configs {
		if route != "default" && len(route) > len(longestMatch) {
			if len(path) >= len(route) && path[:len(route)] == route {
				longestMatch = route
				matchedConfig = config
				foundMatch = true
			}
		}
	}

	if foundMatch {
		return matchedConfig
	}

	// Varsayılan yapılandırma
	if defaultConfig, exists := cm.configs["default"]; exists {
		return defaultConfig
	}

	// Hiçbir şey bulunamazsa çok kısıtlayıcı bir yapılandırma döndür
	return RouteConfig{
		GlobalLimit: 1.0 / 60,
		GlobalBurst: 1,
		UserLimit:   1.0 / 60,
		UserBurst:   1,
	}
}
