package cache_storage

import (
	"github.com/nickham-su/go_memory_cache/cache_value"
	"hash/adler32"
)

// NewFastStorage 创建快速数据存储器
// partitionsCount 分区数量
func NewFastStorage(partitionsCount int) CacheStorage {
	fs := &FastStorage{
		partitions:      make(map[uint32]*MapStorage),
		partitionsCount: uint32(partitionsCount),
	}
	// 初始化分区
	for i := uint32(0); i < fs.partitionsCount; i++ {
		fs.partitions[i] = &MapStorage{storage: make(map[string]*cache_value.CacheValue)}
	}
	return fs
}

// FastStorage 快速数据存储器
// 高频操作时的锁冲突，是影响并发性能的瓶颈
// 解决的思路是分区，使用多个map和锁，降低冲突的几率，从而提升并发性能
// 使用key的哈希分区，映射到对应的MapStorage
type FastStorage struct {
	partitions      map[uint32]*MapStorage // 分区；这里的map是只读的，不用加锁
	partitionsCount uint32                 // 分区数量
}

func (m *FastStorage) getPartition(key string) *MapStorage {
	hash := adler32.New()
	_, _ = hash.Write([]byte(key))
	id := hash.Sum32() % m.partitionsCount
	return m.partitions[id]
}

func (m *FastStorage) Set(key string, value *cache_value.CacheValue) {
	storage := m.getPartition(key)
	storage.Set(key, value)
}

func (m *FastStorage) Get(key string) (value *cache_value.CacheValue, ok bool) {
	storage := m.getPartition(key)
	return storage.Get(key)
}

func (m *FastStorage) Del(key string) {
	storage := m.getPartition(key)
	storage.Del(key)
}

func (m *FastStorage) Exists(key string) bool {
	storage := m.getPartition(key)
	return storage.Exists(key)
}

func (m *FastStorage) Flush() {
	for _, storage := range m.partitions {
		storage.Flush()
	}
}

func (m *FastStorage) Keys() int64 {
	var count int64 = 0
	for _, storage := range m.partitions {
		count += storage.Keys()
	}
	return count
}

func (m *FastStorage) Items() []Item {
	list := make([]Item, 0)
	for _, storage := range m.partitions {
		list = append(list, storage.Items()...)
	}
	return list
}
