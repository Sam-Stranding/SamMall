package admin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"github.com/Sam-Stranding/SamMall/src/utils/tools"
	"github.com/wenlng/go-captcha/v2/slide"
	"go.uber.org/zap"
)

func (s *Service) GetSlideCaptcha(ctx context.Context) (*dto.GetVerifyCaptchaResp, common.Errno) {
	captData, err := s.captcha.Generate()
	if err != nil {
		logger.Error("GetSlideCaptcha Generate Error", zap.Error(err))
		return nil, common.ServerErr.WithErr(err)
	}
	dotData := captData.GetData()
	if dotData == nil {
		logger.Error("GetSlideCaptcha GetData Error")
		return nil, common.ServerErr.WithMsg("GetDat is nil")
	}

	dots, err := json.Marshal(dotData)
	if err != nil {
		logger.Error("GetSlideCaptcha json.Marshal Error", zap.Error(err))
		return nil, common.ServerErr.WithErr(err)
	}

	var mBs64Data, tBs64Data string
	mBs64Data, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		logger.Error("GetSlideCaptcha GetMasterImage.ToBase64 Error", zap.Error(err))
		return nil, common.ServerErr.WithErr(err)
	}
	tBs64Data, err = captData.GetTileImage().ToBase64()
	if err != nil {
		logger.Error("GetSlideCaptcha GetSlaveImage.ToBase64 Error", zap.Error(err))
		return nil, common.ServerErr.WithErr(err)
	}

	key := tools.UUIDHex()
	err = s.verify.SetCaptchaKey(ctx, key, string(dots), time.Minute*2)
	if err != nil {
		logger.Error("GetSlideCaptcha SetCaptchaKey Error", zap.Error(err))
		return nil, common.RedisErr.WithErr(err)
	}

	return &dto.GetVerifyCaptchaResp{
		Expire:         110,
		ImageBs64:      mBs64Data,
		Key:            key,
		TitleHeight:    dotData.Height,
		TitleImageBs64: tBs64Data,
		TitleWidth:     dotData.Width,
		TitleX:         dotData.TileX,
		TitleY:         dotData.TileY,
	}, common.OK

}

func (s *Service) CheckSlideCaptcha(ctx context.Context, req *dto.CheckCaptchaReq) (*dto.CheckCaptchaDtoResp, common.Errno) {
	dots, err := s.verify.GetCaptchaKey(ctx, req.Key)
	if err != nil {
		logger.Error("CheckSlideCaptcha GetCaptchaKey Error", zap.Error(err))
		return nil, common.RedisErr.WithErr(err)
	}
	if dots == "" {
		return nil, common.ParamErr.WithMsg("滑块已过期，请重试")
	}
	dot := slide.Block{}
	err = json.Unmarshal([]byte(dots), &dot)
	if err != nil {
		logger.Error("CheckSlideCaptcha json.Unmarshal Error", zap.Error(err))
		return nil, common.InvalidCaptchaErr
	}
	ok := slide.CheckPoint(int64(req.SlideX), int64(req.SlideY), int64(dot.X), int64(dot.Y), 5)
	if !ok {
		return nil, common.InvalidCaptchaErr
	}
	ticket := tools.UUIDHex()
	err = s.verify.SetCaptchaTicket(ctx, ticket, req.Key, time.Minute*5)
	if err != nil {
		logger.Error("CheckSlideCaptcha SetCaptchaTicket Error", zap.Error(err))
		return nil, common.RedisErr.WithErr(err)
	}
	return &dto.CheckCaptchaDtoResp{
		Expire: 280,
		Ticket: ticket,
	}, common.OK
}
