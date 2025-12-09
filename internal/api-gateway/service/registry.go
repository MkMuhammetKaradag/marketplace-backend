package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"marketplace/internal/api-gateway/config"
)

// ServiceHealth tracks the health status of a service
type ServiceHealth struct {
	Healthy   bool
	LastCheck time.Time
	FailCount int32
	mu        sync.RWMutex
}

// Service represents a backend service
type Service struct {
	Name       string
	BaseURLs   []string
	PathPrefix string
	Health     *ServiceHealth
	Timeout    time.Duration
	nextIndex  uint64
}

// ServiceRegistry manages services
type ServiceRegistry struct {
	services sync.Map // map[string]*Service
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{}
}

func (sr *ServiceRegistry) Register(name string, baseURLs []string, pathPrefix string) error {
	if len(baseURLs) == 0 {
		return fmt.Errorf("at least one BaseURL is required")
	}

	service := &Service{
		Name:       name,
		BaseURLs:   baseURLs,
		PathPrefix: pathPrefix,
		Health: &ServiceHealth{
			Healthy:   true,
			LastCheck: time.Now(),
		},
		Timeout:   config.DefaultTimeout,
		nextIndex: 0,
	}
	sr.services.Store(name, service)
	log.Printf("✅ Service registered: %s -> %v (prefix: %s)", name, baseURLs, pathPrefix)
	return nil
}

func (sr *ServiceRegistry) GetByPath(path string) (*Service, bool) {
	var found *Service
	longestPrefixLen := 0

	sr.services.Range(func(key, value interface{}) bool {
		service := value.(*Service)
		if strings.HasPrefix(path, service.PathPrefix) {
			if len(service.PathPrefix) > longestPrefixLen {
				found = service
				longestPrefixLen = len(service.PathPrefix)
			}
		}
		return true
	})
	return found, found != nil
}

func (sr *ServiceRegistry) List() []*Service {
	var services []*Service
	sr.services.Range(func(key, value interface{}) bool {
		services = append(services, value.(*Service))
		return true
	})
	return services
}

func (s *Service) GetNextBaseURL() (string, bool) {
	s.Health.mu.RLock()
	if !s.Health.Healthy {
		s.Health.mu.RUnlock()
		return "", false
	}
	s.Health.mu.RUnlock()

	index := atomic.AddUint64(&s.nextIndex, 1) - 1
	urlIndex := index % uint64(len(s.BaseURLs))

	return s.BaseURLs[urlIndex], true
}

func (sr *ServiceRegistry) IsHealthy(service *Service) bool {
	service.Health.mu.RLock()
	defer service.Health.mu.RUnlock()
	return service.Health.Healthy
}

func (sr *ServiceRegistry) StartHealthChecks(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			sr.services.Range(func(key, value interface{}) bool {
				service := value.(*Service)
				go sr.checkHealth(service)
				return true
			})
		}
	}()
	log.Printf("🏥 Health check started (interval: %v)", interval)
}

func (sr *ServiceRegistry) checkHealth(service *Service) {
	allHealthy := true

	service.Health.mu.Lock()
	defer service.Health.mu.Unlock()

	service.Health.LastCheck = time.Now()
	client := &http.Client{Timeout: 5 * time.Second}

	for _, baseURL := range service.BaseURLs {
		resp, err := client.Get(baseURL + "/health")
		if err != nil || resp.StatusCode != http.StatusOK {
			atomic.AddInt32(&service.Health.FailCount, 1)
			allHealthy = false
		} else {
			if resp != nil {
				resp.Body.Close()
			}
			atomic.StoreInt32(&service.Health.FailCount, 0)
		}
	}

	if !allHealthy && atomic.LoadInt32(&service.Health.FailCount) >= 3 {
		if service.Health.Healthy {
			log.Printf("🔴 Service DOWN (Circuit Open): %s", service.Name)
			service.Health.Healthy = false
		}
	} else if allHealthy {
		if !service.Health.Healthy {
			log.Printf("🟢 Service UP (Circuit Closed): %s", service.Name)
		}
		service.Health.Healthy = true
	}
}
