package admin

import (
	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/admin"
)

type Service struct {
	adminUser admin.IAdminUser
	user      admin.IAdminUser
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		adminUser: admin.NewAdminUser(adaptor),
		user:      admin.NewAdminUser(adaptor),
	}
}
