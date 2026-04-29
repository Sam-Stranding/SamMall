package do

import "github.com/Sam-Stranding/SamMall/src/common"

type AddRoleReq struct {
	AdminUserID int64
	Name        string
	Desc        string
}

type UpdateRoleReq struct {
	AdminUserID int64
	ID          int64
	Name        string
	Desc        string
	Status      int32
}

type ListRoleReq struct {
	NameKw string
	Status int32
	common.Pager
}
