package token

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/redis"
	"github.com/Sam-Stranding/SamMall/src/adaptor/rpc"
	"github.com/Sam-Stranding/SamMall/src/config"
	"github.com/Sam-Stranding/SamMall/src/consts"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"github.com/gogf/gf/util/gconv"
	"go.uber.org/zap"
)

type AccessToken struct {
	Token    string `json:"token"`
	ExpireIn int64  `json:"expire_in"`
}

type TenantAccessToken struct {
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int64  `json:"expire"`
}

type GetToken func() (*AccessToken, error)

type Service struct {
	conf        *config.Config
	locker      redis.ILocker
	accessToken redis.IAccessToken
	lark        rpc.ILark
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		conf:        adaptor.GetConf(),
		locker:      redis.NewLocker(adaptor),
		accessToken: redis.NewAccessToken(adaptor),
		lark:        rpc.NewLark(adaptor),
	}
}

// 存储access token 的key
func (s *Service) cacheTokenKeyFmt(appCode int32) string {
	return fmt.Sprintf("%s:cacheTokenKey:%d", config.ServerFullName, appCode)
}

// 分布式锁Key
func (s *Service) lockTokenKeyFmt(appCode int32) string {
	return fmt.Sprintf("%s:lockTokenKey:%d", config.ServerFullName, appCode)
}

func (s *Service) UpdateToken(ctx context.Context, getToken GetToken, lockKey, cacheKey string) (*AccessToken, error) {
	locker, err := s.locker.GetLock(ctx, lockKey)
	if err != nil {
		return nil, err
	}
	if locker {
		token, err := getToken()
		if err != nil {
			logger.Error("UpdateToken getToken error", zap.Error(err))
			return nil, err
		}
		err = s.accessToken.Set(ctx, cacheKey, gconv.String(token),
			time.Duration(token.ExpireIn-consts.ExpireTokenDueDuration)*time.Second)
		if err != nil {
			logger.Error("UpdateToken Set error", zap.Error(err))
			return nil, err
		}
		return token, nil
	}
	//等待锁结束
	err = s.locker.AwaitLock(ctx, lockKey, time.Second*2)
	if err != nil {
		logger.Error("UpdateToken AwaitLock error", zap.Error(err))
		return nil, err
	}
	return s.getCache(ctx, cacheKey)
}

func (s *Service) getCache(ctx context.Context, cacheKey string) (*AccessToken, error) {
	tokenJson, expireIn, err := s.accessToken.Get(ctx, cacheKey)
	if err != nil {
		logger.Error("getCache Get error", zap.Error(err))
		return nil, err
	}

	retToken := &AccessToken{}
	err = json.Unmarshal([]byte(tokenJson), &retToken)
	if err != nil {
		logger.Error("getCache json.Unmarshal error", zap.Error(err))
		return nil, err
	}
	retToken.ExpireIn = expireIn
	return retToken, nil
}
