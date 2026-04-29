package admin

import (
	"context"
	"time"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/model"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/query"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/consts"
	"github.com/Sam-Stranding/SamMall/src/service/do"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type IPerm interface {
	PermissionList(ctx context.Context, pager common.Pager) ([]*model.Permission, int64, error)
	MyPermissionList(ctx context.Context, pager common.Pager) ([]*model.Permission, error)
	GetPermNameMap(ctx context.Context, permIds []int64) (map[int64]string, error)

	CreatePermission(ctx context.Context, req *do.AddPerm) (int64, error)
	UpdatePermission(ctx context.Context, req *do.UpdatePermList) error
	DeletePermission(ctx context.Context, req *do.DeletePerm) error
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

func (a *AdminPerm) MyPermissionList(ctx context.Context, pager common.Pager) ([]*model.Permission, error) {
	list, _, err := a.PermissionList(ctx, common.Pager{
		Page:      1,
		Limit:     100,
		UnLimited: true,
	})
	return list, err
}

func (a *AdminPerm) GetPermNameMap(ctx context.Context, permIds []int64) (map[int64]string, error) {
	qs := query.Use(a.db).Permission
	list, err := qs.WithContext(ctx).Where(qs.ID.In(permIds...)).Select(qs.ID, qs.Name).Find()
	if err != nil {
		return nil, err
	}
	retMap := lo.SliceToMap(list, func(item *model.Permission) (int64, string) {
		return item.ID, item.Name
	})
	return retMap, nil
}

func (a *AdminPerm) CreatePermission(ctx context.Context, req *do.AddPerm) (int64, error) {
	qs := query.Use(a.db).Permission
	timeNow := time.Now()
	perm := &model.Permission{
		Code:     req.Code,
		Name:     req.Name,
		PagePath: req.PagePath,
		ParentID: req.ParentID,
		Sort:     req.Sort,
		Type:     req.Type,
		Desc:     req.Desc,
		Status:   consts.IsEnable,
		CreateAt: timeNow,
		UpdateAt: timeNow,
		UpdateBy: req.AdminUserID,
	}
	err := qs.WithContext(ctx).Create(perm)
	if err != nil {
		return 0, err
	}
	return perm.ID, nil
}

func (a *AdminPerm) UpdatePermission(ctx context.Context, req *do.UpdatePermList) error {
	qs := query.Use(a.db).Permission
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range req.List {
			updateMap := map[string]interface{}{
				qs.UpdateBy.ColumnName().String(): item.AdminUserID,
				qs.UpdateAt.ColumnName().String(): time.Now(),
			}
			if item.Code != "" {
				updateMap[qs.Code.ColumnName().String()] = item.Code
			}
			if item.Name != "" {
				updateMap[qs.Name.ColumnName().String()] = item.Name
			}
			if item.PagePath != "" {
				updateMap[qs.PagePath.ColumnName().String()] = item.PagePath
			}
			if item.ParentID != 0 {
				updateMap[qs.ParentID.ColumnName().String()] = item.ParentID
			}
			if item.Sort != 0 {
				updateMap[qs.Sort.ColumnName().String()] = item.Sort
			}
			if item.Type != 0 {
				updateMap[qs.Type.ColumnName().String()] = item.Type
			}
			if item.Desc != "" {
				updateMap[qs.Desc.ColumnName().String()] = item.Desc
			}
			err := tx.Model(&model.Permission{}).Where(qs.ID.Eq(item.ID)).Updates(updateMap).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *AdminPerm) DeletePermission(ctx context.Context, req *do.DeletePerm) error {
	qs := query.Use(a.db).Permission
	_, err := qs.WithContext(ctx).Delete(&model.Permission{
		ID: req.ID,
	})
	return err
}
