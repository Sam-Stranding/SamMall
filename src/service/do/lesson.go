package do

import "github.com/Sam-Stranding/SamMall/src/common"

type AddCategory struct {
	Name     string
	Level    int32
	ParentID int64
	Sort     int32
}

type UpdateCategory struct {
	ID   int64
	Name string
}

type DeleteCategory struct {
	ID int64
}

type UpdateSort struct {
	ID   int64
	Sort int32
}

type UpdateCategorySort []UpdateSort

type CategoryDo struct {
	ID       int64
	Name     string
	Level    int32
	ParentID int64
	Sort     int32
}

type ListCategory struct {
	common.Pager
}
