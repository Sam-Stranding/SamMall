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
	"github.com/go-redis/redis"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type IAdminUser interface {
	CreateUser(ctx context.Context, req *do.CreateUser) (int64, error)
	UpdateUser(ctx context.Context, req *do.UpdateUser) error
	UpdateUserPassword(ctx context.Context, adminUser *do.UpdateUserPassword) error
	UpdateUserLarkOpenID(ctx context.Context, userID int64, openID string) error

	GetUserByMobile(ctx context.Context, mobile string) (*model.AdminUser, error)
	GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error)
	GetUserByLarkOpenID(ctx context.Context, openID string) (*model.AdminUser, error)
	GetOpenIDByMobile(ctx context.Context, mobile string) (string, error)

	ListAdminUser(ctx context.Context, req *do.ListAdminUser) ([]*model.AdminUser, int64, error)
}

type AdminUser struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewAdminUser(adaptor adaptor.IAdaptor) *AdminUser {
	return &AdminUser{
		db:    adaptor.GetDB(),
		redis: adaptor.GetRedis(),
	}
}

func (a *AdminUser) CreateUser(ctx context.Context, req *do.CreateUser) (int64, error) {
	timeNow := time.Now()
	addObject := &model.AdminUser{
		Name:     req.Name,
		NickName: req.NickName,
		Mobile:   req.Mobile,
		Sex:      req.Sex,
		CreateAt: timeNow,
		UpdateAt: timeNow,
		UpdateBy: req.AdminUserID,
		Status:   consts.IsEnable,
		CreateBy: req.AdminUserID,
	}
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(addObject).Error
		if err != nil {
			return err
		}
		userRoles := make([]model.AdminUserRole, 0)
		for _, roleID := range req.RoleIDs {
			userRoles = append(userRoles, model.AdminUserRole{
				AdminUserID: addObject.ID,
				RoleID:      roleID,
				UpdateAt:    timeNow,
				UpdateBy:    req.AdminUserID,
			})
		}
		return tx.CreateInBatches(userRoles, 100).Error
	})
	if err != nil {
		return 0, err
	}
	return addObject.ID, nil
}

func (a *AdminUser) UpdateUser(ctx context.Context, req *do.UpdateUser) error {
	timeNow := time.Now()
	qs := query.Use(a.db).AdminUser
	rqs := query.Use(a.db).AdminUserRole
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Where(qs.ID.Eq(req.ID)).Updates(&model.AdminUser{
			Name:     req.Name,
			NickName: req.NickName,
			Sex:      req.Sex,
			Status:   req.Status,
			UpdateAt: timeNow,
			UpdateBy: req.AdminUserID,
		}).Error
		if err != nil {
			return err
		}
		err = tx.Where(rqs.AdminUserID.Eq(req.ID)).Delete(model.AdminUserRole{}).Error
		if err != nil {
			return err
		}
		userRoles := make([]model.AdminUserRole, 0)
		for _, roleID := range req.RoleIDs {
			userRoles = append(userRoles, model.AdminUserRole{
				AdminUserID: req.ID,
				RoleID:      roleID,
				UpdateAt:    timeNow,
				UpdateBy:    req.AdminUserID,
			})
		}
		return tx.CreateInBatches(userRoles, 100).Error
	})
}

func (a *AdminUser) UpdateUserPassword(ctx context.Context, adminUser *do.UpdateUserPassword) error {
	qs := query.Use(a.db).AdminUser
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(adminUser.ID)).Update(qs.Password, adminUser.Password)
	return err
}

func (a *AdminUser) UpdateUserLarkOpenID(ctx context.Context, userID int64, openID string) error {
	qs := query.Use(a.db).AdminUser
	_, err := qs.WithContext(ctx).Where(qs.ID.Eq(userID)).Update(qs.LarkOpenID, openID)
	return err
}

func (a *AdminUser) GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error) {
	qs := query.Use(a.db).AdminUser
	return qs.WithContext(ctx).Where(qs.ID.Eq(userId)).First()
}

func (a *AdminUser) GetUserByMobile(ctx context.Context, mobile string) (*model.AdminUser, error) {
	qs := query.Use(a.db).AdminUser
	return qs.WithContext(ctx).Where(qs.Mobile.Eq(mobile)).First()
}

func (a *AdminUser) GetUserByLarkOpenID(ctx context.Context, openID string) (*model.AdminUser, error) {
	qs := query.Use(a.db).AdminUser
	return qs.WithContext(ctx).Where(qs.LarkOpenID.Eq(openID)).First()
}

func (a *AdminUser) GetOpenIDByMobile(ctx context.Context, mobile string) (string, error) {
	qs := query.Use(a.db).AdminUser
	var openID string
	err := qs.WithContext(ctx).Select(qs.LarkOpenID).Where(qs.Mobile.Eq(mobile)).Scan(&openID)
	if err != nil {
		return "", err
	}
	return openID, nil
}

func (a *AdminUser) ListAdminUser(ctx context.Context, req *do.ListAdminUser) ([]*model.AdminUser, int64, error) {
	qs := query.Use(a.db).AdminUser
	tx := qs.WithContext(ctx).Where(qs.IsDelete.Neq(consts.IsEnable))
	// 有条件的搜索
	if req.Name != "" {
		tx = tx.Where(qs.Name.Like(tools.GetAllLike(req.Name)))
	}
	if req.Mobile != "" {
		tx = tx.Where(qs.Mobile.Like(tools.GetAllLike(req.Mobile)))
	}
	if req.Status != 0 {
		tx = tx.Where(qs.Status.Eq(req.Status))
	}
	if req.RoleID != 0 {
		rqs := query.Use(a.db).AdminUserRole
		list, err := rqs.WithContext(ctx).Select(rqs.AdminUserID.Distinct()).Where(rqs.RoleID.Eq(req.RoleID)).Find()
		if err != nil {
			return nil, 0, err
		}
		userIds := make([]int64, 0)
		lo.ForEach(list, func(item *model.AdminUserRole, index int) {
			userIds = append(userIds, item.AdminUserID)
		})
		tx = tx.Where(qs.ID.In(userIds...))
	}
	count, err := tx.Count()
	if err != nil {
		return nil, 0, err
	}
	retList, err := tx.Offset(req.GetOffset()).Limit(req.Limit).Order(qs.Status.Desc(), qs.CreateAt.Desc()).Find()

	return retList, count, err
}
