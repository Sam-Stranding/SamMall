package do

import "github.com/Sam-Stranding/SamMall/src/common"

type UserInfo struct{}

type CreateUser struct {
	AdminUserID int64   `json:"admin_user_id"`
	Name        string  `json:"name"`
	NickName    string  `json:"nick_name"`
	Mobile      string  `json:"mobile"`
	Sex         int32   `json:"sex"`
	RoleIDs     []int64 `json:"role_ids"`
}

type UpdateUser struct {
	AdminUserID int64   `json:"admin_user_id"`
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	NickName    string  `json:"nick_name"`
	Sex         int32   `json:"sex"`
	Status      int32   `json:"status"`
	RoleIDs     []int64 `json:"role_ids"`
}

type UpdateUserPassword struct {
	ID       int64  `json:"id"`
	Password string `json:"password"`
}

type ListAdminUser struct {
	common.Pager
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	RoleID int64  `json:"role_id"`
	Status int32  `json:"status"`
}
