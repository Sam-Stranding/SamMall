package do

type LarkUserInfo struct {
	Name            string `json:"name"`             //用户姓名
	EnName          string `json:"en_name"`          //用户英文名
	AvatarURL       string `json:"avatar_url"`       //头像
	AvatarThumb     string `json:"avatar_thumb"`     //头像72*72
	AvatarMiddle    string `json:"avatar_middle"`    //头像240*240
	AvatarBig       string `json:"avatar_big"`       //头像640*640
	OpenID          string `json:"open_id"`          //用户再应用内的唯一标识
	UnionID         string `json:"union_id"`         //用户对ISV的唯一标识
	Email           string `json:"email"`            //邮箱（字段权限要求：获取用户邮箱信息，仅自建应用）
	EnterpriseEmail string `json:"enterprise_email"` //企业邮箱请先确保一再管理后台启用飞书邮箱服务（字段权限要求）
	UserID          string `json:"user_id"`          //用户user_id（用户权限要求：获取用户user_id，仅自建应用）
	Mobile          string `json:"mobile"`           //用户手机号（）
	TenantKey       string `json:"tenant_key"`       //当前企业标识
	EmployeeNo      string `json:"employee_no"`      //用户工号（）
}

type LarkUserInfoResp struct {
	Code int64        `json:"code"`
	Msg  string       `json:"msg"`
	Data LarkUserInfo `json:"data"`
}

type LarkUserAccessToken struct {
	Code        int64  `json:"code"`
	AccessToken string `json:"access_token"`
	ExpireIn    int64  `json:"expire_in"`
	ErrCode     string `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}
