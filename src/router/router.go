package router

import (
	"net/http"

	"github.com/Sam-Stranding/SamMall/src/adaptor"
	"github.com/Sam-Stranding/SamMall/src/api/admin"
	"github.com/Sam-Stranding/SamMall/src/api/customer"
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
	//1.注入鉴权中间件，同时过滤一些不需要鉴权的接口
	app.Use(AuthMiddleware(r.SpanFilter))
	if r.conf.Server.EnablePprof {
		//2.注入pprof
		SetupPprof(app, "/debug/pprof")
	}
	//通过/ping接口检查服务器状态
	app.Any("/ping", r.checkServer())
	root := app.Group(r.rootPath)
	r.route(root)
}

func (r *Router) SpanFilter(ctx *gin.Context) bool {
	return true
}

func (r *Router) AccessRecordFilter(ctx *gin.Context) bool {
	return true
}

func (r *Router) route(root *gin.RouterGroup) {
	adminRoot := root.Group("/admin")
	adminRoot.GET("/user/info", r.admin.GetUserInfo)
}
