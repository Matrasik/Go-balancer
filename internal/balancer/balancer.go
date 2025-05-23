package balancer

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type BackedPool struct {
	BackendsInfo []*BackendServerInfo
	current      int
	mu           sync.Mutex
}

// Смотрм текущий по номеру сервер. Если мертв - идем до первого живого.
func (pool *BackedPool) Next() (addr *url.URL, id int) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	for i := 0; i < len(pool.BackendsInfo); i++ {
		idx := (pool.current + i) % len(pool.BackendsInfo)
		backend := pool.BackendsInfo[idx]

		if backend.IsAlive() {
			pool.current = (idx + 1) % len(pool.BackendsInfo)
			return backend.Address, idx
		}
	}
	return nil, -1
}

// Каждые interval секунд делаем пинг на по урлу /health. Ругаемся, если сервер не доступен
func (pool *BackedPool) HealthCheck(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		pool.mu.Lock()
		for _, elem := range pool.BackendsInfo {
			elem := elem
			go func(b *BackendServerInfo) {
				resp, err := http.Get(b.UrlString + "/health")
				alive := err == nil && resp.StatusCode == 200
				if !alive {
					log.Printf("[Health] Server Id:%d Url:%s unreachable", elem.Id, elem.UrlString)
				}
				b.SetAlive(alive)
			}(elem)
		}
		pool.mu.Unlock()
	}
}

// Тот же HealthChek, но одноразовый
func (pool *BackedPool) InitCheck() {
	for _, elem := range pool.BackendsInfo {
		elem := elem
		go func(b *BackendServerInfo) {
			resp, err := http.Get(b.UrlString + "/health")
			alive := err == nil && resp.StatusCode == 200
			if !alive {
				log.Printf("[Health] Server Id:%d Url:%s unreachable", elem.Id, elem.UrlString)
			}
			b.SetAlive(alive)
		}(elem)
	}
}
