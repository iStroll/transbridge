// cache/cache.go
package cache

import (
	"context"
	"errors"
	"time"
)

// 错误定义
var (
	ErrCacheMiss         = errors.New("cache miss")
	ErrCacheNotAvailable = errors.New("cache not available")
	ErrInvalidTTL        = errors.New("invalid ttl")
)

// Cache 定义缓存接口
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Clear(ctx context.Context) error
	Close(ctx context.Context) error
}

type CacheEntry struct {
	Translation string `json:"translation"`
	Provider    string `json:"provider"`
	APIURL      string `json:"api_url"`
	Model       string `json:"model"`
}
