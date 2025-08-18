package limiter

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*Visitor
	rate     int           // requests
	interval time.Duration // per interval
}

type Visitor struct {
	lastSeen time.Time
	tokens   int
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		interval: interval,
	}
	go rl.cleanupVisitors()
	return rl
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(rl.interval)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.interval {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	now := time.Now()
	if !exists || time.Since(v.lastSeen) > rl.interval {
		rl.visitors[ip] = &Visitor{lastSeen: now, tokens: rl.rate - 1}
		return true
	}
	if v.tokens > 0 {
		v.tokens--
		v.lastSeen = now
		return true
	}
	return false
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if !rl.Allow(ip) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
