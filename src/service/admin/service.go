package admin

import (
	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/admin"
)

type Service struct {
	adminUser admin.IAdminUser
	userInfo  admin.IAdminUser
}

func NewService(adaptor *adaptor.Adaptor) *Service {
	return &Service{
		adminUser: admin.NewAdminUser(adaptor),
	}
}
