package dto

import "github.com/Sam-Stranding/SamMall/src/common"

type AddCategoryReq struct {
	Name     string `json:"name"`
	Level    int32  `json:"level"`
	ParentID int64  `json:"parent_id"`
	Sort     int32  `json:"sort"`
}

type UpdateCategoryReq struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type DeleteCategoryReq struct {
	ID int64 `json:"id"`
}

type UpdateSort struct {
	ID       int64 `json:"id"`
	Sort     int32 `json:"sort"`
	Level    int32 `json:"level"`
	ParentID int64 `json:"parent_id"`
}

type CategoryDto struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Level    int32  `json:"level"`
	ParentID int64  `json:"parent_id"`
	Sort     int32  `json:"sort"`
}

type ListCategoryReq struct {
	common.Pager
}
