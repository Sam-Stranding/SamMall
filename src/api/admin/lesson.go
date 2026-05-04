package admin

import (
	"github.com/Sam-Stranding/SamMall/src/api"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/service/dto"
	"github.com/gin-gonic/gin"
)

func (c *Ctrl) CreateCategory(ctx *gin.Context) {
	req := &dto.AddCategoryReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr)
		return
	}
	resp, errno := c.lesson.CreateCategory(ctx, req)
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) UpdateCategory(ctx *gin.Context) {
	req := &dto.UpdateCategoryReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr)
		return
	}
	errno := c.lesson.UpdateCategory(ctx, req)
	api.WriteResp(ctx, nil, errno)
}

func (c *Ctrl) DeleteCategory(ctx *gin.Context) {
	req := &dto.DeleteCategoryReq{}
	if err := ctx.BindJSON(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr)
		return
	}
	errno := c.lesson.DeleteCategory(ctx, req)
	api.WriteResp(ctx, nil, errno)
}

func (c *Ctrl) CategoryList(ctx *gin.Context) {
	req := &dto.ListCategoryReq{}
	if err := ctx.BindQuery(req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr)
		return
	}
	resp, errno := c.lesson.CategoryList(ctx, req)
	api.WriteResp(ctx, resp, errno)
}

func (c *Ctrl) CategorySorts(ctx *gin.Context) {
	req := []dto.UpdateSort{}
	if err := ctx.BindJSON(&req); err != nil {
		api.WriteResp(ctx, nil, common.ParamErr)
		return
	}
	errno := c.lesson.CategorySorts(ctx, req)
	api.WriteResp(ctx, nil, errno)
}
