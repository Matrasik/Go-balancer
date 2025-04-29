package ratelimiter

import (
	"Golang_balancer/internal/config"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type BucketManager struct {
	buckets    sync.Map
	defaultCfg config.BucketConfig
	clientsCfg map[string]config.BucketConfig
	config.ClientsCfg
}

func NewBucketManager(cfg config.ClientsCfg) *BucketManager {
	clients := make(map[string]config.BucketConfig)
	for _, client := range cfg.Clients {
		clients[client.Addr] = client.BucketConfig
	}
	bm := &BucketManager{
		clientsCfg: clients,
		defaultCfg: config.BucketConfig{
			Capacity: 10, // Дефолтные значения
			Rate:     5,
		},
	}

	for ip, cfg := range clients {
		bm.buckets.Store(ip, &TokenBucket{
			tokens:     cfg.Capacity,
			capacity:   cfg.Capacity,
			rate:       cfg.Rate,
			lastRefill: time.Now().UnixNano(),
		})
	}
	return bm
}

func (bm *BucketManager) GetBucket(ClientIp string) *TokenBucket {
	val, ok := bm.buckets.Load(ClientIp)
	if ok {
		return val.(*TokenBucket)
	}

	cfg, exists := bm.clientsCfg[ClientIp]
	if !exists {
		cfg = bm.defaultCfg
	}

	bucket := &TokenBucket{
		tokens:     cfg.Capacity,
		capacity:   cfg.Capacity,
		rate:       cfg.Rate,
		lastRefill: time.Now().UnixNano(),
	}
	actual, _ := bm.buckets.LoadOrStore(ClientIp, bucket)
	return actual.(*TokenBucket)
}

func (bm *BucketManager) BucketMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bucket := bm.GetBucket(r.RemoteAddr)
		fmt.Printf("This client has %d tokens, %d capacity, %d rate\n", bucket.tokens, bucket.capacity, bucket.rate)
		if !bucket.Allow() {
			http.Error(w, "Limit!!", http.StatusTooManyRequests)
			log.Printf("Client %s reached limit", r.RemoteAddr)
			return
		}
		next.ServeHTTP(w, r)

	})
}
