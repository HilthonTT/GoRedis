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

	return hash % len(m)
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

func (m ShardMap) SAdd(key string, members ...string) int {
	shard := m.getShard(key)

	shard.Lock()
	defer shard.Unlock()

	// Retrieve the existing set or create a new one
	set, ok := shard.data[key].(map[string]struct{})
	if !ok {
		set = make(map[string]struct{})
		shard.data[key] = set
	}

	added := 0
	for _, member := range members {
		if _, exists := set[member]; !exists {
			set[member] = struct{}{}
			added++
		}
	}

	return added
}

func (m ShardMap) SMembers(key string) []string {
	shard := m.getShard(key)

	shard.RLock()
	defer shard.RUnlock()

	set, ok := shard.data[key].(map[string]struct{})
	if !ok {
		return nil
	}

	members := make([]string, 0, len(set))
	for member := range set {
		members = append(members, member)
	}

	return members
}

func (m ShardMap) SRem(key string, members ...string) int {
	shard := m.getShard(key)

	shard.Lock()
	defer shard.Unlock()

	set, ok := shard.data[key].(map[string]struct{})
	if !ok || len(set) == 0 {
		return 0 // nothing to remove
	}

	removed := 0
	for _, member := range members {
		if _, exists := set[member]; exists {
			delete(set, member)
			removed++
		}
	}

	if len(set) == 0 {
		delete(shard.data, key)
	}

	return removed
}
