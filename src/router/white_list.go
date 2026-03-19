package router

var AdminAuthWhiteList = map[string]bool{
	"ping":                          true,
	"metrics":                       true,
	"/admin/v1/user/verify/captcha": true,
}
