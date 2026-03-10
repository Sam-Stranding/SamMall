package admin

import (
	"context"

	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/do"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"go.uber.org/zap"
)

func (s *Service) HelloWorld(ctx context.Context, adminUser *common.AdminUser, req *dto.HelloWorldReq) (*dto.HelloWorldResp, common.Errno) {
	msg, err := s.adminUser.HelloWorld(ctx, &do.HelloWorld{})
	if err != nil {
		logger.Error("HelloWorld HelloWorld err", zap.Error(err), zap.Any("req", req))
		return nil, common.DatabaseErr.WithErr(err)
	}
	return &dto.HelloWorldResp{Hello: msg, World: "world"}, common.OK
}
