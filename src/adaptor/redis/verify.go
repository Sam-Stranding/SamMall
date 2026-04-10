package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/config"
	"github.com/Sam-Stranding/SamMall/src/service/do"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/Sam-Stranding/SamMall/src/utils/http"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type IVerify interface {
	SetCaptchaKey(ctx context.Context, key string, value string, expire time.Duration) error
	GetCaptchaKey(ctx context.Context, key string) (string, error)
	SetCaptchaTicket(ctx context.Context, ticket string, value string, expire time.Duration) error
	GetCaptchaTicket(ctx context.Context, ticket string) (string, error)

	SetAdminUserToken(ctx context.Context, token string, tokenData string, expire time.Duration) error
	GetAdminUserToken(ctx context.Context, token string) (string, error)

	IncrPasswordErr(ctx context.Context, mobile string, expire time.Duration) (int64, error)
	DeletePasswordErr(ctx context.Context, mobile string) error

	GetLarkSmsCode(ctx context.Context, req *dto.GetSmsCodeReq, TenantAccessToken string, UserOpenID string, captcha string) (*do.LarkSmsCodeResp, error)

	SetMobileVerifyCode(ctx context.Context, mobile string, code string, expire time.Duration) error
	GetMobileVerifyCode(ctx context.Context, mobile string) (string, error)
}

type Verify struct {
	redis *redis.Client
}

const (
	larkHost = "https://open.feishu.cn"
	headerCT = "application/json; charset=utf-8"
)

func NewVerify(adaptor adaptor.IAdaptor) *Verify {
	return &Verify{redis: adaptor.GetRedis()}
}

// 封装Key，避免暴露
func fmtVerifyCaptchaKey(key string) string {
	return fmt.Sprintf("%s:captcha:%s", config.ServerFullName, key)
}

func fmtVerifyCaptchaTicket(key string) string {
	return fmt.Sprintf("%s:captcha_ticket:%s", config.ServerFullName, key)
}

func fmtVerifyAdminUserToken(token string) string {
	return fmt.Sprintf("%s:admin_user_token:%s", config.ServerFullName, token)
}

func (v *Verify) SetCaptchaKey(ctx context.Context, key string, value string, expire time.Duration) error {
	redisKey := fmtVerifyCaptchaKey(key)
	return v.redis.Set(redisKey, value, expire).Err()
}

func (v *Verify) GetCaptchaKey(ctx context.Context, key string) (string, error) {
	redisKey := fmtVerifyCaptchaKey(key)
	get, err := v.redis.Get(redisKey).Result()
	if err != nil {
		return "", err
	}
	v.redis.Del(redisKey)
	return get, nil
}

func (v *Verify) SetCaptchaTicket(ctx context.Context, ticket string, value string, expire time.Duration) error {
	redisTicket := fmtVerifyCaptchaTicket(ticket)
	return v.redis.Set(redisTicket, value, expire).Err()
}

func (v *Verify) GetCaptchaTicket(ctx context.Context, ticket string) (string, error) {
	redisTicket := fmtVerifyCaptchaTicket(ticket)
	get, err := v.redis.Get(redisTicket).Result()
	if err != nil {
		return "", err
	}
	//拿了一次就失效，将key删除
	v.redis.Del(redisTicket)
	return get, nil
}

func (v *Verify) SetAdminUserToken(ctx context.Context, token string, tokenData string, expire time.Duration) error {
	redisToken := fmtVerifyAdminUserToken(token)
	return v.redis.Set(redisToken, tokenData, expire).Err()
}

func (v *Verify) GetAdminUserToken(ctx context.Context, token string) (string, error) {
	redisToken := fmtVerifyAdminUserToken(token)
	get, err := v.redis.Get(redisToken).Result()
	if err != nil {
		return "", err
	}
	return get, nil
}

func fmtVerifyMobilePasswordErr(mobile string) string {
	return fmt.Sprintf("%s:admin_user_password_errcount:%s", config.ServerFullName, mobile)
}

func (v *Verify) IncrPasswordErr(ctx context.Context, mobile string, expire time.Duration) (int64, error) {
	redisMobile := fmtVerifyMobilePasswordErr(mobile)
	pipe := v.redis.Pipeline()
	incr, err := pipe.Incr(redisMobile).Result()
	if err != nil {
		return 0, err
	}
	if incr == 1 {
		pipe.Expire(redisMobile, expire)
	}
	_, err = pipe.Exec()
	return incr, err
}

func (v *Verify) DeletePasswordErr(ctx context.Context, mobile string) error {
	redisMobile := fmtVerifyMobilePasswordErr(mobile)
	return v.redis.Del(redisMobile).Err()
}

func (v *Verify) GetLarkSmsCode(ctx context.Context, req *dto.GetSmsCodeReq, TenantAccessToken string, UserOpenID string, captcha string) (*do.LarkSmsCodeResp, error) {
	url := fmt.Sprintf("%s/open-apis/im/v1/messages?receive_id_type=open_id", larkHost)
	headers := map[string]string{
		"Authorization": "Bearer " + TenantAccessToken,
		"Context-Type":  headerCT,
	}
	body := map[string]string{
		"receive_id": UserOpenID,
		"msg_type":   "text",
		"content":    fmt.Sprintf(`{"text":"%s"}`, captcha),
	}
	resp := &do.LarkSmsCodeResp{}
	err := http.Post(ctx, url, headers, body, resp)
	if err != nil {
		logger.Error("GetLarkSmsCode error", zap.Error(err))
		return nil, err
	}
	return resp, nil
}

func fmtMobileVerifyCode(mobile string) string {
	return fmt.Sprintf("%s:mobile_verify:%s", config.ServerFullName, mobile)
}

func (v *Verify) SetMobileVerifyCode(ctx context.Context, mobile string, code string, expire time.Duration) error {
	redisMobile := fmtMobileVerifyCode(mobile)
	return v.redis.Set(redisMobile, code, expire).Err()
}

func (v *Verify) GetMobileVerifyCode(ctx context.Context, mobile string) (string, error) {
	redisMobile := fmtMobileVerifyCode(mobile)
	verify, err := v.redis.Get(redisMobile).Result()
	if err != nil {
		return "", err
	}
	v.redis.Del(redisMobile)
	return verify, nil
}
