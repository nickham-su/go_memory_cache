package cache_storage

import (
	"github.com/nickham-su/go_memory_cache/cache_value"
)

// Item 键值对
type Item struct {
	key   string
	value *cache_value.CacheValue
}

func (i *Item) GetKey() string {
	return i.key
}

func (i Item) GetValue() *cache_value.CacheValue {
	return i.value
}
