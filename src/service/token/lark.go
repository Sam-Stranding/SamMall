package token

import (
	"context"

	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"go.uber.org/zap"
)

func (s *Service) GetLarkUserAccessToken(ctx context.Context, appCode int32, code string, redirectUrl string, scope string, force bool) (*AccessToken, common.Errno) {
	token, err := s.getLarkUserAccessToken(ctx, appCode, code, redirectUrl, scope, force)
	if err != nil {
		logger.Error("GetLarkUserAccessToken get access token failed", zap.Error(err), zap.Int32("appCode", appCode))
		return nil, common.ServerErr.WithErr(err)
	}
	return token, common.OK
}

func (s *Service) getLarkUserAccessToken(ctx context.Context, appCode int32, code string, redirectUrl string, scope string, force bool) (*AccessToken, error) {
	getTokenFunc := func() (*AccessToken, error) {
		token, err := s.lark.GetLarkAccessToken(ctx, appCode, code, redirectUrl, scope)
		if err != nil {
			logger.Error("getLarkUserAccessToken GetUserAccessToken get access token failed", zap.Error(err), zap.Int32("appCode", appCode))
			return nil, common.ServerErr.WithErr(err)
		}
		return &AccessToken{
			ExpireIn: token.ExpireIn,
			Token:    token.AccessToken,
		}, nil
	}
	rpcToken, err := getTokenFunc()
	if err != nil {
		logger.Error("getLarkUserAccessToken getTokenFunc get access token failed", zap.Error(err), zap.Int32("appCode", appCode))
		return nil, common.ServerErr.WithErr(err)
	}
	return rpcToken, common.OK
}
