package ratelimiter

import (
	"Golang_balancer/internal/config"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
	"time"
)

type BucketManager struct {
	buckets    sync.Map
	defaultCfg config.BucketConfig
	clientsCfg map[string]config.BucketConfig
	db         *gorm.DB
	config.ClientsCfg
}

func NewBucketManager(cfg config.ClientsCfg, db *gorm.DB) *BucketManager {
	clients := make(map[string]config.BucketConfig)
	for _, client := range cfg.Clients {
		clients[client.Addr] = client.BucketConfig
	}
	bm := &BucketManager{
		db:         db,
		clientsCfg: clients,
		defaultCfg: config.BucketConfig{
			Capacity: 10,
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
		var dbConfig config.BucketDBConfig
		result := bm.db.Where("ip = ?", ClientIp).First(&dbConfig)
		if result.Error != nil {
			cfg = bm.defaultCfg
			dbConfig = config.BucketDBConfig{
				IP:       ClientIp,
				Capacity: bm.defaultCfg.Capacity,
				Rate:     bm.defaultCfg.Rate,
			}
			if err := bm.db.Create(&dbConfig).Error; err != nil {
				log.Printf("Failed to create bucket config at bd: %v", err)
			}
		} else {
			cfg.Rate = dbConfig.Rate
			cfg.Capacity = dbConfig.Capacity
		}
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
