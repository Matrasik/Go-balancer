package balancer

import (
	"net/url"
	"sync"
)

type BackedPool struct {
	backendsInfo []*BackendServerInfo
	current      int
	mu           sync.Mutex
}

func (pool *BackedPool) Next() *url.URL {
	pool.mu.Lock()
	backend := pool.backendsInfo[pool.current]
	//TODO проверка на живность сервера
	pool.current = (pool.current + 1) % len(pool.backendsInfo)
	pool.mu.Unlock()
	return backend.Address
}
