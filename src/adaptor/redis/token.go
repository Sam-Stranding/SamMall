package redis

import (
	"context"
	"time"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/go-redis/redis"
	"github.com/gogf/gf/util/gconv"
)

type IAccessToken interface {
	Get(ctx context.Context, key string) (string, int64, error)
	Set(ctx context.Context, key string, token string, expireIn time.Duration) error
}

type AccessToken struct {
	redis *redis.Client
}

func NewAccessToken(adaptor adaptor.IAdaptor) *AccessToken {
	return &AccessToken{
		redis: adaptor.GetRedis(),
	}
}

func (a *AccessToken) Get(ctx context.Context, key string) (string, int64, error) {
	tokenJson, err := a.redis.Get(key).Result()
	if err != nil {
		return "", 0, err
	}
	expireIn, err := a.redis.TTL(key).Result()
	if err != nil {
		return "", 0, err
	}
	return tokenJson, gconv.Int64(expireIn.Seconds()), nil
}

func (a *AccessToken) Set(ctx context.Context, key string, token string, expireIn time.Duration) error {
	return a.redis.Set(key, token, expireIn).Err()
}
