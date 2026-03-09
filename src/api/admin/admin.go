package admin

import (
	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/api"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/gin-gonic/gin"
)

type Ctrl struct {
	adaptor *adaptor.Adaptor
}

func NewCtrl(adaptor *adaptor.Adaptor) *Ctrl {
	return &Ctrl{
		adaptor: adaptor,
	}
}

func (c *Ctrl) HelloWorld(ctx *gin.Context) {
	api.WriteResp(ctx, "hello world", common.OK)
}
