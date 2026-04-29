package admin

import (
	"github.com/Sam-Stranding/SamMall/src/api"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/gin-gonic/gin"
)

func (c *Ctrl) PermissionList(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	resp, errno := c.perm.PermissionList(ctx.Request.Context())
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) MyPermissionList(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	resp, errno := c.perm.MyPermissionList(ctx.Request.Context(), user)
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) CreatePermission(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.AddPermissionReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr)
		return
	}
	permID, errno := c.perm.CreatePermission(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, map[string]interface{}{
		"id": permID,
	}, errno)
}

func (c *Ctrl) UpdatePermission(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.UpdatePermissionReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr)
		return
	}
	errno := c.perm.UpdatePermissions(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}

func (c *Ctrl) DeletePermission(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.DeletePermissionReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr)
		return
	}
	errno := c.perm.DeletePermission(ctx.Request.Context(), user, req)
	api.WriteResp(ctx, nil, errno)
}
