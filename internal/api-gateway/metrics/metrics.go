package metrics

import (
	"fmt"
	"sync"
	"time"
)

type Metrics struct {
	mu                  sync.RWMutex
	totalRequests       int64
	successRequests     int64
	failedRequests      int64
	rateLimitedRequests map[string]int64
	requestsByService   map[string]int64
	requestsByPath      map[string]int64
	lastReset           time.Time
}

func NewMetrics() *Metrics {
	return &Metrics{
		rateLimitedRequests: make(map[string]int64),
		requestsByService:   make(map[string]int64),
		requestsByPath:      make(map[string]int64),
		lastReset:           time.Now(),
	}
}

func (m *Metrics) IncrementTotal() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalRequests++
}

func (m *Metrics) IncrementSuccess() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.successRequests++
}

func (m *Metrics) IncrementFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failedRequests++
}

func (m *Metrics) IncrementRateLimit(limitType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rateLimitedRequests[limitType]++
}

func (m *Metrics) IncrementService(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestsByService[service]++
}

func (m *Metrics) IncrementPath(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestsByPath[path]++
}

func (m *Metrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	uptime := time.Since(m.lastReset)

	rps := 0.0
	if uptime.Seconds() > 0 {
		rps = float64(m.totalRequests) / uptime.Seconds()
	}

	successRate := 0.0
	if m.totalRequests > 0 {
		successRate = float64(m.successRequests) / float64(m.totalRequests) * 100
	}
	return map[string]interface{}{
		"total_requests":        m.totalRequests,
		"success_requests":      m.successRequests,
		"failed_requests":       m.failedRequests,
		"success_rate":          fmt.Sprintf("%.2f%%", successRate),
		"rate_limited_requests": m.rateLimitedRequests,
		"requests_by_service":   m.requestsByService,
		"requests_by_path":      m.requestsByPath,
		"uptime_seconds":        uptime.Seconds(),
		"requests_per_second":   rps,
	}
}
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalRequests = 0
	m.successRequests = 0
	m.failedRequests = 0
	m.rateLimitedRequests = make(map[string]int64)
	m.requestsByService = make(map[string]int64)
	m.requestsByPath = make(map[string]int64)
	m.lastReset = time.Now()
}