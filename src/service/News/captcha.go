package News

import (
	"context"

	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"go.uber.org/zap"
)

func (s *Service) GetLarkSmsCode(ctx context.Context, req *dto.GetSmsCodeReq, TenantAccessToken string, UserOpenID string) (*SmsCode, error) {
	captcha, errno := s.getLarkSmsCode(ctx, req, TenantAccessToken, UserOpenID)
	if errno.NotOk() {
		logger.Error("GetLarkSmsCode getLarkSmsCode error", zap.Error(errno), zap.String("OpenID", UserOpenID))
		return nil, errno
	}
	return captcha, nil
}

func (s *Service) getLarkSmsCode(ctx context.Context, req *dto.GetSmsCodeReq, TenantAccessToken string, UserOpenID string) (*SmsCode, common.Errno) {
	var captchaCode string
	captchaCode = GenerateRandomNumber()
	getNewsFunc := func() (*SmsCode, error) {
		smsCode, err := s.verify.GetLarkSmsCode(ctx, req, TenantAccessToken, UserOpenID, captchaCode)
		if err != nil {
			logger.Error("getLarkSmsCode GetLarkSmsCode error", zap.Error(err))
			return &SmsCode{
				ErrCode: smsCode.Code,
				ErrMsg:  "获取验证码失败",
			}, nil
		}
		return &SmsCode{
			Code: captchaCode,
			Msg:  "验证码发送成功",
		}, nil
	}
	redisSmsCode, err := getNewsFunc()
	if err != nil {
		logger.Error("getLarkSmsCode getNewsFunc error", zap.Error(err))
	}
	return redisSmsCode, common.OK
}
