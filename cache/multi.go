package cache

import (
	"context"
	"time"
)

type MultiCache struct {
	caches []Cache
}

func NewMultiCache(caches []Cache) *MultiCache {
	return &MultiCache{
		caches: caches,
	}
}

// Close 实现 Cache 接口
func (m *MultiCache) Close(ctx context.Context) error {
	var lastErr error
	for _, cache := range m.caches {
		if err := cache.Close(ctx); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Get 从多级缓存中获取数据
func (m *MultiCache) Get(ctx context.Context, key string) (string, error) {
	var lastErr error
	for i, cache := range m.caches {
		value, err := cache.Get(ctx, key)
		if err == nil {
			// 找到数据后，更新之前的缓存层
			for j := 0; j < i; j++ {
				m.caches[j].Set(ctx, key, value, time.Hour*24)
			}
			return value, nil
		}
		lastErr = err
	}
	return "", lastErr
}

// Set 设置多级缓存的数据
func (m *MultiCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	var lastErr error
	for _, cache := range m.caches {
		if err := cache.Set(ctx, key, value, ttl); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Clear 清除所有缓存
func (m *MultiCache) Clear(ctx context.Context) error {
	var lastErr error
	for _, cache := range m.caches {
		if err := cache.Clear(ctx); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
