package ratelimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	tokens     int64
	capacity   int64
	rate       int64
	lastRefill int64 //Решил не использовать тикер😣
	mu         sync.Mutex
}

func (tb *TokenBucket) refill(now int64) {
	difference := now - tb.lastRefill
	addTokens := (difference * tb.rate) / 1e9
	if addTokens+tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
		return
	}
	tb.tokens += addTokens
	tb.lastRefill = now
	return
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	now := time.Now().UnixNano()
	tb.refill(now)
	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}
