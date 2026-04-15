package admin

import (
	"context"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/model"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/query"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/consts"
	"gorm.io/gorm"
)

type IPerm interface {
	PermissionList(ctx context.Context, pager common.Pager) ([]*model.Permission, int64, error)
	MyPermissionList(ctx context.Context, pager common.Pager) ([]*model.Permission, int64, error)
}

type AdminPerm struct {
	db *gorm.DB
}

func NewAdminPerm(adaptor adaptor.IAdaptor) *AdminPerm {
	return &AdminPerm{
		db: adaptor.GetDB(),
	}
}

func (a *AdminPerm) PermissionList(ctx context.Context, pager common.Pager) ([]*model.Permission, int64, error) {
	qs := query.Use(a.db).Permission
	return qs.WithContext(ctx).Where(qs.Status.Eq(consts.IsEnable)).FindByPage(pager.GetOffset(), pager.Limit)

}

func (a *AdminPerm) MyPermissionList(ctx context.Context, pager common.Pager) ([]*model.Permission, int64, error) {
	return a.PermissionList(ctx, pager)
}
