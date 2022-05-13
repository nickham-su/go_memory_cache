package cache_storage

import (
	"github.com/nickham-su/go_memory_cache/cache_value"
)

// CacheStorage 仅负责数据存储（协程安全）
type CacheStorage interface {
	// Set 设置⼀个缓存项
	Set(key string, value *cache_value.CacheValue)
	// Get 获取⼀个值
	Get(key string) (value *cache_value.CacheValue, ok bool)
	// Del 删除⼀个值
	Del(key string)
	// Exists 检测⼀个值 是否存在
	Exists(key string) bool
	// Flush 清空所有值
	Flush()
	// Keys 返回所有的key 多少
	Keys() int64
	// Items 获取所有键值对
	Items() []Item
}
