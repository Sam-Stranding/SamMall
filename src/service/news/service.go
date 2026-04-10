package news

import (
	"context"
	"crypto/rand"
	"math/big"
	"strconv"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/redis"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/config"
	"github.com/Sam-Stranding/SamMall/src/consts"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"go.uber.org/zap"
)

type SmsCode struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	ErrCode int64  `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

type Service struct {
	conf   *config.Config
	verify redis.IVerify
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		conf:   adaptor.GetConf(),
		verify: redis.NewVerify(adaptor),
	}
}

// GenerateRandomNumber 随机生成验证码
func GenerateRandomNumber() string {
	n, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(n.Int64()) + 1000)
}

// StoreMobileVerify 存储手机号+验证码，过期时间：5分钟
func (s *Service) StoreMobileVerify(ctx context.Context, mobile string, code string) error {
	return s.verify.SetMobileVerifyCode(ctx, mobile, code, consts.MobileVerifyExpire)
}

// VerifyMobileVerifyCode 验证手机号+验证码是否正确
func (s *Service) VerifyMobileVerifyCode(ctx context.Context, mobile string, inputVerify string) (bool, error) {
	storeVerify, err := s.verify.GetMobileVerifyCode(ctx, mobile)
	if err != nil {
		logger.Error("VerifyMobileVerifyCode GetMobileVerifyCode error", zap.Error(err), zap.String("mobile", mobile))
		return false, err
	}
	if inputVerify == storeVerify {
		return true, nil
	} else {
		logger.Error("VerifyMobileVerifyCode verify error", zap.String("mobile", mobile), zap.String("inputVerify", inputVerify), zap.String("storeVerify", storeVerify))
		return false, common.MobileVerifyIncorrectErr
	}
}
