package go_memory_cache

import (
	"errors"
	"github.com/nickham-su/go_memory_cache/cache_storage"
	"github.com/nickham-su/go_memory_cache/cache_value"
	"github.com/nickham-su/go_memory_cache/memory_manager"
	"time"
)

// CleaningTimeInterval 清理过期任务的时间间隔
const CleaningTimeInterval = time.Minute

// ErrMemoryFulled 内存达到上限错误
var ErrMemoryFulled = errors.New("memory is full")

// NewMemoryCache 创建内存缓存
// fastMode false是基于单个map存储；true使用分区降低锁冲突，提升性能
func NewMemoryCache(fastMode bool) Cache {
	var storage cache_storage.CacheStorage
	if fastMode {
		storage = cache_storage.NewFastStorage(128) // 快速数据存储器
	} else {
		storage = cache_storage.NewMapStorage() // 基于map的储存数据
	}

	memoryCache := &MemoryCache{
		memoryManager: memory_manager.NewMemoryManager(),
		storage:       storage,
	}

	// 定时处理过期数据
	go func() {
		for range time.Tick(CleaningTimeInterval) {
			memoryCache.clearingExpired()
		}
	}()
	return memoryCache
}

// MemoryCache 内存缓存
type MemoryCache struct {
	memoryManager *memory_manager.MemoryManager // 内存管理器
	storage       cache_storage.CacheStorage    // 数据存储器
}

// SetMaxMemory size 是⼀个字符串。⽀持以下参数: 1KB，100KB，1MB，2MB，1GB 等
func (m *MemoryCache) SetMaxMemory(size string) bool {
	return m.memoryManager.SetMaxMemory(size)
}

// Set 设置⼀个缓存项，并且在expire时间之后过期
// expire留空或为0，则永不过期
func (m *MemoryCache) Set(key string, val interface{}, expire ...time.Duration) error {
	if m.memoryManager.IsFull() {
		return ErrMemoryFulled
	}
	if len(expire) == 0 {
		m.storage.Set(key, cache_value.NewCacheValue(val, cache_value.NeverExpire))
	} else {
		m.storage.Set(key, cache_value.NewCacheValue(val, expire[0]))
	}
	return nil
}

// Get 获取⼀个值
func (m *MemoryCache) Get(key string) (interface{}, bool) {
	value, ok := m.storage.Get(key)
	if !ok || value.Expired() {
		return nil, false
	}
	return value.GetValue(), true
}

// Del 删除⼀个值
// 成功返回true;如果返回false，则没有对应的key
func (m *MemoryCache) Del(key string) bool {
	m.storage.Del(key)
	return true
}

// Exists 检测⼀个值 是否存在
func (m *MemoryCache) Exists(key string) bool {
	return m.storage.Exists(key)
}

// Flush 清空所有值
func (m *MemoryCache) Flush() bool {
	m.storage.Flush()
	return true
}

// Keys 返回所有的key的数量
func (m *MemoryCache) Keys() int64 {
	m.clearingExpired()
	return m.storage.Keys()
}

// clearingExpired 清理过期数据
func (m *MemoryCache) clearingExpired() {
	items := m.storage.Items()
	deleteKeys := make([]string, 0)
	for _, item := range items {
		if item.GetValue().Expired() {
			deleteKeys = append(deleteKeys, item.GetKey())
		}
	}
	for _, key := range deleteKeys {
		m.storage.Del(key)
	}
}
