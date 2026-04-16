package admin

import (
	"context"

	"github.com/Sam-Stranding/SamMall/src/api"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/gin-gonic/gin"
)

func (c *Ctrl) GetAdminUserByToken(ctx context.Context, token string) (*common.AdminUser, error) {
	adminUser, errno := c.user.GetAdminUserByToken(ctx, token)
	if errno.NotOk() {
		return nil, errno
	}
	return adminUser, nil
}

func (c *Ctrl) AdminUserList(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.ListAdminUserReq{}
	resp, errno := c.user.AdminUserList(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) GetUserInfo(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	resp, errno := c.user.GetUserInfo(ctx.Request.Context(), user)
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) CreateUser(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.CreateUserReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithMsg(err.Error()))
	}
	userID, errno := c.user.CreateUser(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, map[string]int64{
		"id": userID,
	}, errno)
}

func (c *Ctrl) UpdateUser(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.UpdateUserReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithMsg(err.Error()))
	}
	errno := c.user.UpdateUser(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

func (c *Ctrl) LarkBind(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.LarkQrCodeBindReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr.WithMsg(err.Error()))
	}
	errno := c.user.LarkBind(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

func (c *Ctrl) LarkUnbind(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	errno := c.user.LarkUnbind(ctx.Request.Context(), user)
	api.WriteResp(ctx, nil, errno)
}

func (c *Ctrl) AdminUserLogout(ctx *gin.Context) {
	adminUser := api.GetAdminTokenFromCtx(ctx)
	errno := c.user.AdminUserLogout(ctx.Request.Context(), adminUser)
	api.WriteResp(ctx, nil, errno)
}
