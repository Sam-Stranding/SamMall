package api

import (
	"net/http"

	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/consts"
	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	ErrMsg string `json:"err_msg"`
	Data   any    `json:"data"`
}

func WriteResp(ctx *gin.Context, data any, err common.Errno) {
	ctx.JSON(http.StatusOK, Resp{
		Code:   err.Code,
		Msg:    err.Msg,
		ErrMsg: err.ErrMsg,
		Data:   data,
	})
}

func GetTokenFromCtx(ctx *gin.Context) *common.User {
	user, exits := ctx.Get(consts.CustomerUserKey)
	if !exits {
		return nil
	}
	return user.(*common.User)
}

func GetAdminTokenFromCtx(ctx *gin.Context) *common.AdminUser {
	user, exits := ctx.Get(consts.AdminUserKey)
	if !exits {
		return nil
	}
	return user.(*common.AdminUser)
}
