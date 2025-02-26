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

// RedisCacheProvider 是缓存提供者的结构体定义
type RedisCacheProvider struct {
	client *redis.Client
}

// NewRedisCacheProvider 创建一个新的缓存提供者实例
func NewRedisCacheProvider(addr string, password string) *RedisCacheProvider {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // 没有密码则留空
	})

	// 测试连接是否成功
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("无法连接到Redis: %v", err))
	}

	return &RedisCacheProvider{
		client: rdb,
	}
}

// Set 设置缓存值
func (cp *RedisCacheProvider) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return cp.client.SetNX(ctx, key, value, expiration).Err()
}

// Get 获取缓存值
func (cp *RedisCacheProvider) Get(ctx context.Context, key string) (string, error) {
	return cp.client.Get(ctx, key).Result()
}

// Delete 删除缓存值
func (cp *RedisCacheProvider) Delete(ctx context.Context, key string) error {
	return cp.client.Del(ctx, key).Err()
}
