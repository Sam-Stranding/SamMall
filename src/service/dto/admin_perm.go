package dto

import "github.com/Sam-Stranding/SamMall/src/common"

type PermissionDto struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`      // 权限编码
	Type     int32  `json:"type"`      // 1:菜单  2：操作
	Name     string `json:"name"`      // 权限名称
	PagePath string `json:"page_path"` // 菜单路径
	ParentID int64  `json:"parent_id"` // 父级权限ID
	Status   int32  `json:"status"`    // 1：正常 -1：禁用
	Sort     int32  `json:"sort"`
	Desc     string `json:"desc"` // 权限描述
}
type PermissionListResp struct {
	common.Pager
	Total int64            `json:"total"`
	List  []*PermissionDto `json:"list"`
}

type AddPermissionReq struct {
	Code     string `json:"code"`
	Type     int32  `json:"type"`
	Desc     string `json:"desc"`
	Name     string `json:"name"`
	PagePath string `json:"page_path"`
	ParentID int64  `json:"parent_id"`
	Sort     int32  `json:"sort"`
}

type UpdatePermDto struct {
	ID int64 `json:"id"`
	AddPermissionReq
}

type UpdatePermissionReq struct {
	List []UpdatePermDto `json:"list"`
}

type DeletePermissionReq struct {
	ID int64 `json:"id"`
}
