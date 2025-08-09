package expiration

import (
	"goredis-server/internal/cache"
	"sync"
	"time"
)

var (
	mu          sync.RWMutex
	Expirations = make(map[string]time.Time)
	stopCh      = make(chan struct{})
)

func SetExpiration(key string, ttl time.Duration) {
	mu.Lock()
	Expirations[key] = time.Now().Add(ttl)
	mu.Unlock()
}

func RemoveExpiration(key string) {
	mu.Lock()
	delete(Expirations, key)
	mu.Unlock()
}

func StartExpirationCleaner(db *cache.ShardMap) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cleanupExpired(db)
			case <-stopCh:
				return
			}
		}
	}()
}

func StopExpirationCleaner() {
	close(stopCh)
}

func cleanupExpired(db *cache.ShardMap) {
	now := time.Now()

	var expired []string
	mu.RLock()
	for key, exp := range Expirations {
		if now.After(exp) {
			expired = append(expired, key)
		}
	}

	mu.RUnlock()

	if len(expired) > 0 {
		mu.Lock()
		for _, key := range expired {
			db.Delete(key)
			delete(Expirations, key)
		}
		mu.Unlock()
	}
}
