package admin

import (
	"github.com/Sam-Stranding/SamMall/src/api"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/gin-gonic/gin"
)

func (c *Ctrl) GetUserInfo(ctx *gin.Context) {
	user := api.GetAdminTokenFromCtx(ctx)
	if user == nil {
		api.WriteResp(ctx, nil, common.AuthErr)
		return
	}
	resp, errno := c.user.GetUserInfo(ctx.Request.Context(), &common.AdminUser{})
	api.WriteResp(ctx, resp, errno)
}
