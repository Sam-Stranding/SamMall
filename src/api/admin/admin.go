package admin

import (
	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/service/admin"
	"github.com/Sam-Stranding/SamMall/src/service/perm"
	"github.com/Sam-Stranding/SamMall/src/service/role"
)

type Ctrl struct {
	adaptor adaptor.IAdaptor
	user    *admin.Service
	perm    *perm.Service
	role    *role.Service
}

func NewCtrl(adaptor adaptor.IAdaptor) *Ctrl {
	return &Ctrl{
		adaptor: adaptor,
		user:    admin.NewService(adaptor),
		perm:    perm.NewService(adaptor),
		role:    role.NewService(adaptor),
	}
}
