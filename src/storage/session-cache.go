package storage

import (
	"sync"
	"time"
)

type CacheItem struct {
	timestamp time.Time //last accessed time
	value     interface{}
}

// SessionCache defines a safe in mem cache for session data
type SessionCache struct {
	sync.RWMutex
	keys map[string]*CacheItem
}

func NewSessionCache(expiry time.Duration) *SessionCache {
	sc := &SessionCache{
		keys: make(map[string]*CacheItem),
	}

	go func(sc *SessionCache) {
		for {
			sc.scanAndRemoveExpiredTokens(expiry)
			time.Sleep(time.Second)
		}
	}(sc)

	return sc
}

func (sc *SessionCache) scanAndRemoveExpiredTokens(expiry time.Duration) {
	sc.Lock()
	defer sc.Unlock()

	expired := []string{}
	//scan
	for k, v := range sc.keys {
		d := time.Now().UTC().Sub(v.timestamp)
		if d >= expiry { //expired
			expired = append(expired, k)
		}
	}
	//delete
	for _, k := range expired {
		delete(sc.keys, k)
	}
}

// Set a token
func (sc *SessionCache) Set(key string, value interface{}) {
	sc.Lock()
	defer sc.Unlock()

	ci := &CacheItem{
		timestamp: time.Now().UTC(),
		value:     value,
	}

	sc.keys[key] = ci
}

// Get a SessionKey by token
func (sc *SessionCache) Get(key string) (interface{}, bool) {
	sc.RLock()
	defer sc.RUnlock()

	if sk, ok := sc.keys[key]; ok {
		sk.timestamp = time.Now().UTC()
		return sk.value, true
	}
	return nil, false
}
