package role

import (
	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/admin"
)

type Service struct {
	adminRole admin.IRole
	adminPerm admin.IPerm
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		adminRole: admin.NewAdminRole(adaptor),
		adminPerm: admin.NewAdminPerm(adaptor),
	}
}
