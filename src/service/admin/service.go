package admin

import (
	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/adaptor/redis"
	"github.com/Sam-Stranding/SamMall/src/adaptor/repo/admin"
	"github.com/Sam-Stranding/SamMall/src/adaptor/rpc"
	"github.com/Sam-Stranding/SamMall/src/service/News"
	"github.com/Sam-Stranding/SamMall/src/service/token"
	"github.com/Sam-Stranding/SamMall/src/utils/captcha"
	"github.com/wenlng/go-captcha/v2/slide"
)

type Service struct {
	adminUser admin.IAdminUser
	user      admin.IAdminUser
	captcha   slide.Captcha
	verify    redis.IVerify
	token     *token.Service
	lark      rpc.ILark
	news      *News.Service
}

func NewService(adaptor adaptor.IAdaptor) *Service {
	return &Service{
		adminUser: admin.NewAdminUser(adaptor),
		user:      admin.NewAdminUser(adaptor),
		verify:    redis.NewVerify(adaptor),
		captcha:   captcha.NewSlideCaptcha(),
		token:     token.NewService(adaptor),
		lark:      rpc.NewLark(adaptor),
		news:      News.NewService(adaptor),
	}
}
