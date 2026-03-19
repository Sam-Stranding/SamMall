package admin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"github.com/Sam-Stranding/SamMall/src/utils/tools"
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
