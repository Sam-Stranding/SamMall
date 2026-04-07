package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/config"
	"github.com/Sam-Stranding/SamMall/src/service/do"
	"github.com/Sam-Stranding/SamMall/src/utils/http"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"go.uber.org/zap"
)

const (
	larkHost = "https://open.feishu.cn"
	headerCT = "application/json; charset=utf-8"
)

type ILark interface {
	GetLarkUserInfo(ctx context.Context, userAccessToken string) (*do.LarkUserInfo, error)
	GetLarkAccessToken(ctx context.Context,
		appCode int32, code string,
		redirectUrl string, scope string) (*do.LarkUserAccessToken, error)
}

type Lark struct {
	conf *config.Config
}

func NewLark(adaptor adaptor.IAdaptor) *Lark {
	return &Lark{
		conf: adaptor.GetConf(),
	}
}

func (l *Lark) GetLarkUserInfo(ctx context.Context, userAccessToken string) (*do.LarkUserInfo, error) {
	url := fmt.Sprintf("%s/open-apis/authen/v1/user_info", larkHost)
	headers := map[string]string{
		"Content-Type":  headerCT,
		"Authorization": "Bearer " + userAccessToken,
	}
	err, respBody := http.Get(ctx, url, headers, nil)
	if err != nil {
		logger.Error("GetLarkUserInfo error", zap.Error(err))
		return nil, err
	}
	resp, err := parseLarkUserInfoResponse(respBody)
	if err != nil {
		logger.Error("GetLarkUserInfo parse response error", zap.Error(err))
		return nil, err
	}
	return resp, nil
}

func parseLarkUserInfoResponse(respBody []byte) (*do.LarkUserInfo, error) {
	resp := &do.LarkUserInfoResp{}
	if err := json.Unmarshal(respBody, resp); err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, fmt.Errorf("lark get user info failed: code=%d, msg=%s", resp.Code, resp.Msg)
	}
	return &resp.Data, nil
}

func (l *Lark) GetLarkAccessToken(ctx context.Context,
	appCode int32, code string,
	redirectUrl string, scope string) (*do.LarkUserAccessToken, error) {
	url := "https://open.feishu.cn/open-apis/authen/v2/oauth/token"

	body := map[string]interface{}{
		"grant_type":    "authorization_code",
		"client_id":     l.conf.AppConf[appCode].AppID,
		"client_secret": l.conf.AppConf[appCode].AppSecret,
		"code":          code,
		"redirect_uri":  redirectUrl,
		"scope":         scope,
	}
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json;charset=utf-8",
	}
	resp := &do.LarkUserAccessToken{}
	err := http.Post(ctx, url, headers, body, resp)
	if err != nil {
		logger.Error("GetLarkAccessToken error", zap.Error(err))
		return nil, err
	}
	return resp, nil
}
