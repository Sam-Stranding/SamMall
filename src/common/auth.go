package common

type AdminUser struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	NickName   string `json:"nick_name"`
	Sex        int32  `json:"sex"`
	Status     int32  `json:"status"`
	Mobile     string `json:"mobile"`
	LarkOpenID string `json:"lark_open_id"`
}

type User struct {
	UserID   int64  `json:"user_id"`
	NickName string `json:"nick_name"`
}
