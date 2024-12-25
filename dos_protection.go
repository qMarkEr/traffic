package main

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

const requrst_limit = 0

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu       sync.Mutex
	ipLimits = make(map[string]*ipLimiter)
)

// Создаем лимитер для IP
func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	// Если IP уже есть в мапе, обновляем lastSeen
	if limiter, exists := ipLimits[ip]; exists {
		limiter.lastSeen = time.Now()
		return limiter.limiter
	}

	// Создаем новый лимитер: 10 запросов в секунду
	limiter := rate.NewLimiter(requrst_limit, 20)
	ipLimits[ip] = &ipLimiter{limiter: limiter, lastSeen: time.Now()}

	// Очищаем старые записи
	go cleanupOldEntries()
	return limiter
}

// Удаляем IP-адреса, которые не были активны более минуты
func cleanupOldEntries() {
	mu.Lock()
	defer mu.Unlock()

	for ip, limiter := range ipLimits {
		if time.Since(limiter.lastSeen) > time.Minute {
			delete(ipLimits, ip)
		}
	}
}

// Middleware для ограничения запросов
func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		limiter := getLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
