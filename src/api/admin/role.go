package admin

import (
	"github.com/Sam-Stranding/SamMall/src/api"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/gin-gonic/gin"
)

func (c *Ctrl) RoleList(ctx *gin.Context) {
	req := &dto.ListRoleReq{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr)
		return
	}
	resp, errno := c.role.RoleList(ctx, req)
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) MyRoles(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	resp, errno := c.role.GetMyRoles(ctx, user)
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) AddRole(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.AddRoleReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr)
		return
	}
	id, errno := c.role.CreateRole(ctx, user, req)
	api.WriteResp(ctx, map[string]interface{}{
		"id": id,
	}, errno)
}
func (c *Ctrl) UpdateRole(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.UpdateRoleReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr)
		return
	}
	errno := c.role.UpdateRole(ctx, user, req)
	api.WriteResp(ctx, nil, errno)
}

func (c *Ctrl) SetRolePerms(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	req := &dto.SetRolePermReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr)
		return
	}
	errno := c.role.SetRolePerms(ctx, user, req)
	api.WriteResp(ctx, nil, errno)
}
