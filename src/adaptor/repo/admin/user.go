package admin

import (
	"context"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/service/do"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type IAdminUser interface {
	HelloWorld(ctx context.Context, req *do.HelloWorld) (string, error)
}

type AdminUser struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewAdminUser(adaptor *adaptor.Adaptor) *AdminUser {
	return &AdminUser{
		db:    adaptor.DbClient,
		redis: adaptor.Redis,
	}
}

func (a *AdminUser) HelloWorld(ctx context.Context, req *do.HelloWorld) (string, error) {
	return "hello world", nil
}
