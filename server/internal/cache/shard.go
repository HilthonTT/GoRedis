package cache

import (
	"crypto/sha1"
	"sync"
)

type Shard struct {
	sync.RWMutex
	data map[string]any
}

type ShardMap []*Shard

func NewShardMap(n int) ShardMap {
	shards := make([]*Shard, n)

	for i := range n {
		shards[i] = &Shard{
			data: make(map[string]any),
		}
	}

	return shards
}

func (m ShardMap) getShardIndex(key string) int {
	checksum := sha1.Sum([]byte(key))
	hash := int(checksum[0])

	return hash & len(m)
}

func (m ShardMap) getShard(key string) *Shard {
	i := m.getShardIndex(key)
	return m[i]
}

func (m ShardMap) Get(key string) (any, bool) {
	shard := m.getShard(key)

	shard.RLock()
	defer shard.RUnlock()

	val, ok := shard.data[key]
	return val, ok
}

func (m ShardMap) Set(key string, val any) {
	shard := m.getShard(key)

	shard.Lock()
	defer shard.Unlock()

	shard.data[key] = val
}

func (m ShardMap) Delete(key string) {
	shard := m.getShard(key)

	shard.Lock()
	defer shard.Unlock()

	delete(shard.data, key)
}

func (m ShardMap) Keys() []string {
	keys := make([]string, 0)

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	wg.Add(len(m))

	for _, shard := range m {
		go func(s *Shard) {
			s.RLock()

			for k := range s.data {
				mutex.Lock()
				keys = append(keys, k)
				mutex.Unlock()
			}

			s.RUnlock()
			wg.Done()
		}(shard)
	}

	wg.Wait()

	return keys
}

func (m ShardMap) KeyValues() []KeyValuePair {
	kv := make([]KeyValuePair, 0)

	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	wg.Add(len(m))

	for _, shard := range m {
		go func(s *Shard) {
			s.RLock()

			for k := range s.data {
				mutex.Lock()

				val, ok := s.data[k]
				if !ok {
					continue
				}
				kv = append(kv, KeyValuePair{
					Key:   k,
					Value: val,
				})
				mutex.Unlock()
			}

			s.RUnlock()
			wg.Done()
		}(shard)
	}

	wg.Wait()

	return kv
}
