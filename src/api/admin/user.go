package admin

import (
	"github.com/Sam-Stranding/SamMall/src/api"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/gin-gonic/gin"
)

func (c *Ctrl) GetUserInfo(ctx *gin.Context) {
	//token common.AdminUser
	resp, errno := c.user.GetUserInfo(ctx.Request.Context(), &common.AdminUser{})
	api.WriteResp(ctx, resp, errno)
}
