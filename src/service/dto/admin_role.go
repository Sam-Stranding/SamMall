package dto

import (
	"github.com/Sam-Stranding/SamMall/src/common"
)

type AddRoleReq struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type UpdateRoleReq struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Status int32  `json:"status"`
}

type ListRoleReq struct {
	common.Pager
	NameKw string `url:"name_kw"`
	Status int32  `url:"status"`
}

type RoleDto struct {
	ID       int64           `json:"id"`
	Name     string          `json:"name"`
	Desc     string          `json:"desc"`
	Status   int32           `json:"status"`
	Perms    []common.IDName `json:"perms"`
	CreateAt int64           `json:"create_at"`
	UpdateAt int64           `json:"update_at"`
}

type ListRoleResp struct {
	common.Pager
	List  []*RoleDto `json:"list"`
	Total int64      `json:"total"`
}

type SetRolePermReq struct {
	RoleID  int64   `json:"role_id"`
	PermIDs []int64 `json:"perm_ids"`
}
