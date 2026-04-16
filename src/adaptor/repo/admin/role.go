package admin

import (
	"context"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/model"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/query"
	"github.com/go-redis/redis"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type IAdminRole interface {
	GetRoleByUserIds(ctx context.Context, userIds []int64) (map[int64][]*model.AdminUserRole, error)
	GetRoleByIds(ctx context.Context, roleIds []int64) (map[int64]*model.Role, error)
}

type AdminRole struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewAdminRole(adaptor adaptor.IAdaptor) *AdminRole {
	return &AdminRole{
		db:    adaptor.GetDB(),
		redis: adaptor.GetRedis(),
	}
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
