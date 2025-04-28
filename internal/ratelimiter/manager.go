package ratelimiter

import (
	"Golang_balancer/internal/config"
	"log"
	"net/http"
	"sync"
	"time"
)

type BucketManager struct {
	buckets sync.Map
	Config  *config.BucketConfig
}

func (bm *BucketManager) GetBucket(UserId string) *TokenBucket {
	val, ok := bm.buckets.Load(UserId)
	if !ok {
		bucket := &TokenBucket{
			tokens:     bm.Config.Capacity,
			capacity:   bm.Config.Capacity,
			rate:       bm.Config.Rate,
			lastRefill: time.Now().UnixNano(),
		}
		actual, loaded := bm.buckets.LoadOrStore(UserId, bucket)
		if loaded {
			val = actual
		} else {
			val = bucket
		}
	}
	return val.(*TokenBucket)
}

func (bm *BucketManager) BucketMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bucket := bm.GetBucket(r.RemoteAddr)
		if !bucket.Allow() {
			http.Error(w, "Limit!!", http.StatusTooManyRequests)
			log.Printf("Client %s reached limit", r.RemoteAddr)
			return
		}
		next.ServeHTTP(w, r)

	})
}
