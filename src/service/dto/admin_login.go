package dto

import (
	"fmt"

	"github.com/Sam-Stranding/SamMall/src/utils/tools"
)

type GetVerifyCaptchaReq struct {
	Once string `form:"once"`
	Time int64  `form:"ts"`
	Sign string `form:"sign"` // 秘钥固定加密： md5(once+Sam2026+ts) 转小写
}

func (r *GetVerifyCaptchaReq) CheckSign() bool {
	return r.Sign == tools.Sha256Hash(fmt.Sprintf("%s%s%d", r.Once, "Sam2026", r.Time))
}

type GetVerifyCaptchaResp struct {
	Key            string `json:"key"`
	ImageBs64      string `json:"image_base64"`       // 包含“data:image/jpeg;base64
	TitleImageBs64 string `json:"title_image_base64"` // 滑块图片，包含“data:image/jpeg;base64
	TitleHeight    int    `json:"title_height"`       // 滑块图片高
	TitleWidth     int    `json:"title_width"`        // 滑块图片宽
	TitleX         int    `json:"title_x"`            // 滑块图的x坐标
	TitleY         int    `json:"title_y"`            // 滑块图的y坐标
	Expire         int64  `json:"expire"`             // 过期时间
}

type CheckCaptchaReq struct {
	Key    string `json:"key"`
	SlideX int    `json:"slide_x"`
	SlideY int    `json:"slide_y"`
}

type CheckCaptchaDtoResp struct {
	Ticket string `json:"ticket"` //票据
	Expire int64  `json:"expire"` //过期时间
}

type GetSmsCodeReq struct {
	Mobile  string `json:"mobile"`
	Ticket  string `json:"ticket"`
	AppCode int32  `json:"app_code"`
}

type GetSmsCodeResp struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	ErrCode int64  `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

type MobileVerifyLoginReq struct {
	Mobile  string `json:"mobile"`
	Captcha string `json:"captcha"`
}

type MobilePasswordLoginReq struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Ticket   string `json:"ticket"`
}

type LoginResp struct {
	Token string       `json:"token"`
	User  AdminUserDto `json:"user"`
}

type LarkQrCodeLoginReq struct {
	AppCode     int32  `json:"app_code"`
	Code        string `json:"code"`
	RedirectUri string `json:"redirect_uri"`
}
