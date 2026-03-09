package admin

import (
	"context"

	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
)

func (s *Service) HelloWorld(ctx context.Context, adminUser *common.AdminUser, req *dto.HelloWorldReq) (*dto.HelloWorldResp, common.Errno) {
	return &dto.HelloWorldResp{Hello: "hello", World: "world"}, common.OK
}
