package router

import (
	"context"
	"net/http"
	"strings"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/api/admin"
	"github.com/Sam-Stranding/SamMall/src/api/customer"
	"github.com/Sam-Stranding/SamMall/src/common"
	"github.com/Sam-Stranding/SamMall/src/config"
	"github.com/gin-gonic/gin"
)

type IRouter interface {
	Register(engine *gin.Engine)            //注册
	SpanFilter(r *gin.Context) bool         //跨度过滤器
	AccessRecordFilter(r *gin.Context) bool //访问记录过滤器
}

type Router struct {
	FullPPROF bool
	rootPath  string
	conf      *config.Config
	checkFunc func() error
	admin     *admin.Ctrl
	customer  *customer.Ctrl
}

func NewRouter(adaptor adaptor.IAdaptor, conf *config.Config, checkFunc func() error) *Router {
	return &Router{
		FullPPROF: conf.Server.EnablePprof,
		rootPath:  "/api/mall",
		conf:      conf,
		checkFunc: checkFunc,
		admin:     admin.NewCtrl(adaptor),
		customer:  customer.NewCtrl(adaptor),
	}
}

// 检查服务器状态
func (r *Router) checkServer() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		err := r.checkFunc()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}

func (r *Router) Register(app *gin.Engine) {
	if r.conf.Server.EnablePprof {
		//2.注入pprof
		SetupPprof(app, "/debug/pprof")
	}
	//通过/ping接口检查服务器状态
	app.Any("/ping", r.checkServer())
	root := app.Group(r.rootPath)
	r.route(root)
}

// SpanFilter 过滤器
func (r *Router) SpanFilter(ctx *gin.Context) bool {
	path := strings.Replace(ctx.Request.URL.Path, r.rootPath, "", 1)
	if _, ok := AdminAuthWhiteList[path]; ok {
		return false
	}
	return true
}

// AccessRecordFilter 日志过滤器
func (r *Router) AccessRecordFilter(ctx *gin.Context) bool {
	return true
}

func (r *Router) route(root *gin.RouterGroup) {
	r.customerRoute(root)
	r.adminRoute(root)
}

func (r *Router) customerRoute(root *gin.RouterGroup) {
	//注入鉴权中间件
	cstRoot := root.Group("/customer", AuthMiddleware(r.SpanFilter, func(ctx context.Context, token string) (*common.User, error) {
		return &common.User{}, nil
	}))
	cstRoot.GET("/user/info", r.admin.GetUserInfo)
}

func (r *Router) adminRoute(root *gin.RouterGroup) {
	//注入鉴权中间件
	adminRoot := root.Group("/admin", AdminAuthMiddleware(r.SpanFilter, func(ctx context.Context, token string) (*common.AdminUser, error) {
		return &common.AdminUser{
			UserID: 1,
			Name:   "admin",
		}, nil
	}))
	//登录无鉴权：添加白名单
	adminRoot.GET("v1/user/verify/captcha", r.admin.GetSmsCodeCaptcha)
	adminRoot.POST("v1/user/verify/captcha/check", r.admin.CheckSmsCodeCaptcha)
	adminRoot.GET("/v1/user/info", r.admin.GetUserInfo)
	adminRoot.POST("/v1/user/create", r.admin.CreateUser)
	adminRoot.POST("/v1/user/update", r.admin.UpdateUser)
	adminRoot.POST("/v1/user/update_status", r.admin.UpdateUserStatus)
}
