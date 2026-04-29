package do

type AddPerm struct {
	AdminUserID int64
	Code        string
	Type        int32
	Name        string
	PagePath    string
	ParentID    int64
	Sort        int32
	Desc        string
}

type UpdatePerm struct {
	ID int64
	AddPerm
}

type UpdatePermList struct {
	List []UpdatePerm
}

type DeletePerm struct {
	ID int64
}
