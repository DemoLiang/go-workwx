package workwx

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	AccessTokenKey = "access_token"
)

// CacheProvider 是缓存提供者的结构体定义
type CacheProvider struct {
	client *redis.Client
}

// NewCacheProvider 创建一个新的缓存提供者实例
func NewCacheProvider(addr string, password string, db int) *CacheProvider {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // 没有密码则留空
		DB:       db,       // 使用默认DB
	})

	// 测试连接是否成功
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("无法连接到Redis: %v", err))
	}

	return &CacheProvider{
		client: rdb,
	}
}

// Set 设置缓存值
func (cp *CacheProvider) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return cp.client.SetNX(ctx, key, value, expiration).Err()
}

// Get 获取缓存值
func (cp *CacheProvider) Get(ctx context.Context, key string) (string, error) {
	return cp.client.Get(ctx, key).Result()
}

// Delete 删除缓存值
func (cp *CacheProvider) Delete(ctx context.Context, key string) error {
	return cp.client.Del(ctx, key).Err()
}
