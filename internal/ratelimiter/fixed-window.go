package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowRateLimiter struct {
	sync.RWMutex
	clients map[string]int
	limit   int
	window  time.Duration
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {

	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}

func (fwrl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	fwrl.Lock()
	count, exist := fwrl.clients[ip]
	fwrl.Unlock()

	if !exist || count < fwrl.limit {
		fwrl.Lock()

		if !exist {
			go fwrl.resetCount(ip)
		}
		fwrl.clients[ip]++
		fwrl.Unlock()
		return true, 0
	}

	return false, fwrl.window
}

func (fwrl *FixedWindowRateLimiter) resetCount(ip string) {
	time.Sleep(fwrl.window)

	fwrl.Lock()
	delete(fwrl.clients, ip)
	fwrl.Unlock()
}
