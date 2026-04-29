package admin

import (
	"context"
	"time"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/model"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/query"
	"github.com/Sam-Stranding/SamMall/src/consts"
	"github.com/Sam-Stranding/SamMall/src/service/do"
	"github.com/Sam-Stranding/SamMall/src/utils/tools"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type IRole interface {
	CreateRole(ctx context.Context, req *do.AddRoleReq) (int64, error)
	UpdateRole(ctx context.Context, req *do.UpdateRoleReq) error
	SetRolePerms(ctx context.Context, roleID int64, permIDs []int64, userId int64) error
	GetRolePerms(ctx context.Context, roleID []int64) (map[int64][]int64, error)
	ListRoles(ctx context.Context, req *do.ListRoleReq) ([]*model.Role, int64, error)
	GetRolesByUserID(ctx context.Context, userID int64) ([]*model.AdminUserRole, error)

	GetRoleByUserIds(ctx context.Context, userIds []int64) (map[int64][]*model.AdminUserRole, error)
	GetRoleByIds(ctx context.Context, roleIds []int64) (map[int64]*model.Role, error)
}

type AdminRole struct {
	db *gorm.DB
}

func NewAdminRole(adaptor adaptor.IAdaptor) *AdminRole {
	return &AdminRole{
		db: adaptor.GetDB(),
	}
}

func (a *AdminRole) CreateRole(ctx context.Context, req *do.AddRoleReq) (int64, error) {
	qs := query.Use(a.db).Role
	timeNow := time.Now()
	role := &model.Role{
		Name:     req.Name,
		Desc:     req.Desc,
		Status:   consts.IsEnable,
		CreateBy: req.AdminUserID,
		CreateAt: timeNow,
		UpdateBy: req.AdminUserID,
		UpdateAt: timeNow,
	}
	err := qs.WithContext(ctx).Create(role)
	if err != nil {
		return 0, err
	}
	return role.ID, nil
}

func (a *AdminRole) UpdateRole(ctx context.Context, req *do.UpdateRoleReq) error {
	qs := query.Use(a.db).Role
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updateMap := map[string]interface{}{
			qs.UpdateBy.ColumnName().String(): req.AdminUserID,
			qs.UpdateAt.ColumnName().String(): time.Now(),
		}
		if req.Name != "" {
			updateMap[qs.Name.ColumnName().String()] = req.Name
		}
		if req.Desc != "" {
			updateMap[qs.Desc.ColumnName().String()] = req.Desc
		}
		if req.Status != 0 {
			updateMap[qs.Status.ColumnName().String()] = req.Status
		}
		err := tx.Model(&model.Role{}).Where(qs.ID.Eq(req.ID)).Updates(updateMap).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *AdminRole) SetRolePerms(ctx context.Context, roleID int64, permIDs []int64, userId int64) error {
	qs := query.Use(a.db).RolePermission
	timeNow := time.Now()
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.RolePermission{}).Where(qs.RoleID.Eq(roleID)).Delete(&model.RolePermission{}).Error
		if err != nil {
			return err
		}
		rolePerms := make([]*model.RolePermission, 0)
		for _, permIds := range permIDs {
			rolePerms = append(rolePerms, &model.RolePermission{
				RoleID:       roleID,
				PermissionID: permIds,
				CreateAt:     timeNow,
				UpdateAt:     timeNow,
				CreateBy:     userId,
				UpdateBy:     userId,
			})
		}
		return tx.CreateInBatches(rolePerms, 100).Error
	})
}

func (a *AdminRole) GetRolePerms(ctx context.Context, roleID []int64) (map[int64][]int64, error) {
	qs := query.Use(a.db).RolePermission
	list, err := qs.WithContext(ctx).Where(qs.RoleID.In(roleID...)).Find()
	if err != nil {
		return nil, err
	}
	rolePermMaps := make(map[int64][]int64)
	lo.ForEach(list, func(item *model.RolePermission, index int) {
		rolePermMaps[item.RoleID] = append(rolePermMaps[item.RoleID], item.PermissionID)
	})
	return rolePermMaps, nil
}

func (a *AdminRole) ListRoles(ctx context.Context, req *do.ListRoleReq) ([]*model.Role, int64, error) {
	qs := query.Use(a.db).Role
	tx := qs.WithContext(ctx)
	if req.NameKw != "" {
		tx = tx.Where(qs.Name.Like(tools.GetAllLike(req.NameKw)))
	}
	if req.Status != 0 {
		tx = tx.Where(qs.Status.Eq(req.Status))
	}
	return tx.Order(qs.Status.Desc(), qs.CreateAt.Desc()).FindByPage(req.GetOffset(), req.Limit)
}

func (a *AdminRole) GetRolesByUserID(ctx context.Context, userID int64) ([]*model.AdminUserRole, error) {
	qs := query.Use(a.db).AdminUserRole
	return qs.WithContext(ctx).Where(qs.AdminUserID.Eq(userID)).Find()
}

func (a *AdminRole) GetRoleByUserIds(ctx context.Context, userIds []int64) (map[int64][]*model.AdminUserRole, error) {
	qs := query.Use(a.db).AdminUserRole
	retList, err := qs.WithContext(ctx).Where(qs.AdminUserID.In(userIds...)).Find()
	if err != nil {
		return nil, err
	}
	return lo.GroupBy(retList, func(item *model.AdminUserRole) int64 {
		return item.AdminUserID
	}), nil
}

func (a *AdminRole) GetRoleByIds(ctx context.Context, roleIds []int64) (map[int64]*model.Role, error) {
	qs := query.Use(a.db).Role
	retList, err := qs.WithContext(ctx).Where(qs.ID.In(roleIds...)).Find()
	if err != nil {
		return nil, err
	}
	//将Slice转化为Map，
	//一个场景，如果是Slice，知道roleId，要获取它对应的roleName，需要遍历，时间复杂度为O(N * M)
	//如果是Map，知道roleId，要获取它对应的roleName，只需要roleMap[user.RoleId]，时间复杂度为O(N)
	return lo.SliceToMap(retList, func(item *model.Role) (int64, *model.Role) { return item.ID, item }), nil
}
