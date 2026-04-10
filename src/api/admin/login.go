package admin

import (
	"github.com/Sam-Stranding/SamMall/src/api"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/gin-gonic/gin"
)

func (c *Ctrl) GetSmsCodeCaptcha(ctx *gin.Context) {
	req := &dto.GetVerifyCaptchaReq{}
	if err := ctx.BindQuery(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}
	resp, errno := c.user.GetSlideCaptcha(ctx.Request.Context())
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) CheckSmsCodeCaptcha(ctx *gin.Context) {
	req := &dto.CheckCaptchaReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}
	resp, errno := c.user.CheckSlideCaptcha(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) GetSmsCode(ctx *gin.Context) {
	req := &dto.GetSmsCodeReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}
	resp, errno := c.user.GetSmsCode(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) MobilePasswordLogin(ctx *gin.Context) {
	req := &dto.MobilePasswordLoginReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}
	resp, errno := c.user.MobilePasswordLogin(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) MobileVerifyLogin(ctx *gin.Context) {
	req := &dto.MobileVerifyLoginReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}
	resp, errno := c.user.MobileVerifyLogin(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) LarkQrCodeLogin(ctx *gin.Context) {
	req := dto.LarkQrCodeLoginReq{}
	if err := ctx.BindJSON(&req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithErr(err))
		return
	}
	resp, errno := c.user.LarkQrCodeLogin(ctx.Request.Context(), req)
	api.WriteResp(ctx, resp, errno)
}
