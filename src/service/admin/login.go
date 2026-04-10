package admin

import (
	"context"
	"errors"

	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/consts"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"github.com/Sam-Stranding/SamMall/src/utils/tools"
	"github.com/go-redis/redis"
	"github.com/gogf/gf/util/gconv"
	"go.uber.org/zap"
)

func (s *Service) processToken(ctx context.Context, token string, adminUser *dto.AdminUserDto) error {
	err := s.verify.SetAdminUserToken(ctx, token, gconv.String(adminUser), consts.AdminUserTokenExpire)
	if err != nil {
		logger.Error("SetAdminUserToken Error", zap.Error(err), zap.String("mobile", adminUser.Mobile))
		return common.DatabaseErr.WithErr(err)
	}
	return nil
}

func (s *Service) MobilePasswordLogin(ctx context.Context, req *dto.MobilePasswordLoginReq) (*dto.LoginResp, common.Errno) {
	//ticket校验，手机号+密码登录也要用滑块校验，所以会有ticket
	_, err := s.verify.GetCaptchaTicket(ctx, req.Ticket)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, common.InvalidCaptchaErr.WithMsg("验证码已过期")
		}
		logger.Error("MobilePasswordLogin GetCaptchaTicket Error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, common.ParamErr.WithErr(err)
	}

	adminUser, err := s.user.GetUserByMobile(ctx, req.Mobile)
	if err != nil {
		logger.Error("MobilePasswordLogin GetUserByMobile Error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, common.DatabaseErr.WithErr(err)
	}
	if adminUser == nil || adminUser.Status != consts.IsEnable {
		return nil, common.InvalidPasswordErr
	}
	//进行用户密码累计次数校验
	errCount, err := s.verify.IncrPasswordErr(ctx, req.Mobile, consts.PasswordErrExpire)
	if err != nil {
		logger.Error("MobilePasswordLogin IncrPasswordErr Error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, common.RedisErr.WithErr(err)
	}
	if errCount >= consts.PasswordErrMaxCount {
		//限制密码错误次数，例如10分钟之内只能错误3次
		return nil, common.PasswordErrLimit
	}
	if adminUser.Password != req.Password {
		return nil, common.InvalidPasswordErr
	}
	_ = s.verify.DeletePasswordErr(ctx, req.Mobile)

	adminUserDto := dto.AdminUserDto{
		UserID:     adminUser.ID,
		Name:       adminUser.Name,
		NickName:   adminUser.NickName,
		Sex:        adminUser.Sex,
		Status:     adminUser.Status,
		Mobile:     adminUser.Mobile,
		LarkOpenID: adminUser.LarkOpenID,
		UpdateAt:   adminUser.UpdateAt.UnixMilli(),
		CreateAt:   adminUser.CreateAt.UnixMilli(),
	}
	tokenUuid := tools.UUIDHex()
	err = s.processToken(ctx, tokenUuid, &adminUserDto)
	if err != nil {
		logger.Error("MobilePasswordLogin processToken Error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, common.RedisErr.WithErr(err)
	}
	return &dto.LoginResp{
		Token: tokenUuid,
		User:  adminUserDto,
	}, common.OK
}

func (s *Service) GetSmsCode(ctx context.Context, req *dto.GetSmsCodeReq) (*dto.GetSmsCodeResp, common.Errno) {
	TenantAccessToken, errno := s.token.GetTenantAccessToken(ctx, req.AppCode)
	if errno.NotOk() {
		logger.Error("GetSmsCode GetTenantAccessToken Error", zap.Error(errno), zap.String("mobile", req.Mobile))
		return nil, common.ServerErr.WithErr(errno)
	}
	UserOpenID, err := s.user.GetOpenIDByMobile(ctx, req.Mobile)
	if err != nil {
		logger.Error("GetSmsCode GetOpenIDByMobile Error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, common.ServerErr.WithErr(err)
	}
	SmsCode, err := s.news.GetLarkSmsCode(ctx, req, TenantAccessToken.TenantAccessToken, UserOpenID)
	if err != nil {
		logger.Error("GetSmsCode GetLarkSmsCode Error", zap.Error(err), zap.String("mobile", req.Mobile))
		return &dto.GetSmsCodeResp{
			ErrCode: SmsCode.ErrCode,
			ErrMsg:  SmsCode.ErrMsg,
		}, common.ServerErr.WithErr(err)
	}
	return &dto.GetSmsCodeResp{
		Code: SmsCode.Code,
		Msg:  SmsCode.Msg,
	}, common.OK
}

func (s *Service) MobileVerifyLogin(ctx context.Context, req *dto.MobileVerifyLoginReq) (*dto.MobileVerifyLoginResp, common.Errno) {
	//验证码校验
	_, err := s.news.VerifyMobileVerifyCode(ctx, req.Mobile, req.Captcha)
	if err != nil {
		logger.Error("MobileVerifyLogin VerifyMobileVerifyCode Error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, common.ServerErr.WithErr(err)
	}

	//获取用户信息
	adminUser, err := s.user.GetUserByMobile(ctx, req.Mobile)
	if err != nil {
		logger.Error("MobileVerifyLogin GetUserByMobile Error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, common.DatabaseErr.WithErr(err)
	}
	if adminUser == nil || adminUser.Status != consts.IsEnable {
		return nil, common.InvalidPasswordErr
	}
	adminUserDto := dto.AdminUserDto{
		UserID:     adminUser.ID,
		Name:       adminUser.Name,
		NickName:   adminUser.NickName,
		Sex:        adminUser.Sex,
		Status:     adminUser.Status,
		Mobile:     adminUser.Mobile,
		LarkOpenID: adminUser.LarkOpenID,
		UpdateAt:   adminUser.UpdateAt.UnixMilli(),
		CreateAt:   adminUser.CreateAt.UnixMilli(),
	}
	tokenUuid := tools.UUIDHex()
	err = s.processToken(ctx, tokenUuid, &adminUserDto)
	if err != nil {
		logger.Error("MobilePasswordLogin processToken Error", zap.Error(err), zap.String("mobile", req.Mobile))
		return nil, common.RedisErr.WithErr(err)
	}
	return &dto.MobileVerifyLoginResp{
		Token: tokenUuid,
		User:  adminUserDto,
	}, common.OK
}

func (s *Service) LarkQrCodeLogin(ctx context.Context, req dto.LarkQrCodeLoginReq) (interface{}, common.Errno) {
	accessToken, errno := s.token.GetLarkUserAccessToken(ctx, req.AppCode, req.Code, req.RedirectUri, "", false)
	if errno.NotOk() {
		logger.Error("LarkQrCodeLogin GetLarkUserAccessToken Error", zap.Error(errno), zap.Any("req", req))
		return nil, common.ServerErr.WithErr(errno)
	}
	larkUserInfo, err := s.lark.GetLarkUserInfo(ctx, accessToken.Token)
	if err != nil {
		logger.Error("LarkQrCodeLogin GetLarkUserInfo Error", zap.Error(err), zap.Any("req", req))
		return nil, common.ServerErr.WithErr(err)
	}
	adminUser, err := s.user.GetUserByLarkOpenID(ctx, larkUserInfo.OpenID)
	if err != nil {
		logger.Error("LarkQrCodeLogin GetUserByLarkOpenID Error", zap.Error(err), zap.Any("req", req))
		return nil, common.DatabaseErr.WithErr(err)
	}
	if adminUser == nil || adminUser.Status != consts.IsEnable {
		return nil, common.AdminUSerNotFoundErr
	}
	adminUserDto := dto.AdminUserDto{
		UserID:     adminUser.ID,
		Name:       adminUser.Name,
		NickName:   adminUser.NickName,
		Sex:        adminUser.Sex,
		Status:     adminUser.Status,
		Mobile:     adminUser.Mobile,
		LarkOpenID: adminUser.LarkOpenID,
		UpdateAt:   adminUser.UpdateAt.UnixMilli(),
		CreateAt:   adminUser.CreateAt.UnixMilli(),
	}
	tokenUuid := tools.UUIDHex()
	err = s.processToken(ctx, tokenUuid, &adminUserDto)
	if err != nil {
		logger.Error("LarkQrCodeLogin processToken Error", zap.Error(err), zap.Any("req", req))
		return nil, common.RedisErr.WithErr(err)
	}
	return &dto.LoginResp{
		Token: tokenUuid,
		User:  adminUserDto,
	}, common.OK
}
