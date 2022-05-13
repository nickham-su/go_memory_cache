package cache_storage

import (
	"github.com/nickham-su/go_memory_cache/cache_value"
	"sync"
)

func NewMapStorage() CacheStorage {
	return &MapStorage{
		storage: make(map[string]*cache_value.CacheValue),
	}
}

// MapStorage 基于map的储存数据
// 也可以使用sync.Map实现，不过猜测这个题的意图是处理协程安全，使用封装好的库好像并不合适:)
type MapStorage struct {
	storage map[string]*cache_value.CacheValue
	lock    sync.RWMutex
}

func (m *MapStorage) Set(key string, value *cache_value.CacheValue) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.storage[key] = value
}

func (m *MapStorage) Get(key string) (value *cache_value.CacheValue, ok bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	value, ok = m.storage[key]
	return
}

func (m *MapStorage) Del(key string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.storage, key)
}

func (m *MapStorage) Exists(key string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	_, ok := m.storage[key]
	return ok
}

func (m *MapStorage) Flush() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.storage = make(map[string]*cache_value.CacheValue)
}

func (m *MapStorage) Keys() int64 {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return int64(len(m.storage))
}

func (m *MapStorage) Items() []Item {
	m.lock.RLock()
	defer m.lock.RUnlock()
	list := make([]Item, 0)
	for k, v := range m.storage {
		list = append(list, Item{
			key:   k,
			value: v,
		})
	}
	return list
}
