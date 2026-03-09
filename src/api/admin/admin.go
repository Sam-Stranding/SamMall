package admin

import (
	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/api"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/admin"
	"github.com/gin-gonic/gin"
)

type Ctrl struct {
	adaptor *adaptor.Adaptor
	hello   *admin.Service
}

func NewCtrl(adaptor *adaptor.Adaptor) *Ctrl {
	return &Ctrl{
		adaptor: adaptor,
		hello:   admin.NewService(adaptor),
	}
}

func (c *Ctrl) HelloWorld(ctx *gin.Context) {
	resp, errno := c.hello.HelloWorld(ctx.Request.Context(), &common.AdminUser{}, nil)
	api.WriteResp(ctx, resp, errno)
}
