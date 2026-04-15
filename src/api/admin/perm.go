package admin

import (
	"github.com/Sam-Stranding/SamMall/src/api"
	"github.com/Sam-Stranding/SamMall/src/common"
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
