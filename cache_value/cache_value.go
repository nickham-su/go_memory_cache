package cache_value

import (
	"time"
)

const (
	// NeverExpire 永不过期
	NeverExpire time.Duration = 0
)

// NewCacheValue 创建CacheValue
func NewCacheValue(value interface{}, expire time.Duration) *CacheValue {
	var expireTime time.Time
	if expire != NeverExpire {
		expireTime = time.Now().Add(expire)
	}
	return &CacheValue{
		value:       value,
		expireTime:  expireTime,
		neverExpire: expire == NeverExpire,
	}
}

// CacheValue 缓存的值
type CacheValue struct {
	value       interface{} // 值
	expireTime  time.Time   // 过期时间
	neverExpire bool        // 永不过期
}

// GetValue 获取未过期的值
func (c *CacheValue) GetValue() interface{} {
	return c.value
}

// Expired 是否过期
func (c *CacheValue) Expired() bool {
	return !c.neverExpire && time.Now().After(c.expireTime)
}
