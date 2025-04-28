package balancer

import (
	"net/url"
	"sync"
)

type BackendServerInfo struct {
	Id        int      `json:"id"`
	Address   *url.URL `json:"-"`
	UrlString string   `json:"address"`
	isAlive   bool
	mu        sync.Mutex
}

func (b *BackendServerInfo) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.isAlive = alive
}

func (b *BackendServerInfo) IsAlive() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.isAlive
}
