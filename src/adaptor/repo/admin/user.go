package admin

import (
	"context"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/model"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/query"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type IAdminUser interface {
	GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error)
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

func (a *AdminUser) GetUserInfo(ctx context.Context, userId int64) (*model.AdminUser, error) {
	qs := query.Use(a.db).AdminUser
	return qs.WithContext(ctx).Where(qs.ID.Eq(userId)).First()
}
