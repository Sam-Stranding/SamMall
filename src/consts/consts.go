package consts

import "time"

const (
	UserTokenKey    = "token"
	AdminTokenKey   = "token"
	CustomerUserKey = "user_key"
	AdminUserKey    = "admin_user_key"
)

const (
	ExpireTokenDueDuration = 200
)

const (
	AdminUserTokenExpire = time.Hour * 24
	PasswordErrExpire    = time.Minute * 10
	PasswordErrMaxCount  = 3
)

const (
	IsEnable  = 1
	ISDisable = -1
)
