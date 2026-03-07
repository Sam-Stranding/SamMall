package router

import (
	"github.com/Sam-Stranding/SamMall/src/config"
	"github.com/gin-gonic/gin"
)

type IRouter interface {
	Register(engine *gin.Engine)
	SpanFilter(r *gin.Context) bool
	AccessRecordFilter(r *gin.Context) bool
}

type Router struct {
	FullPPROF bool
	rootPath  string
	conf      *config.Config
	checkFunc func() error
}

func NewRouter(conf *config.Config, checkFunc func() error) *Router {
	return &Router{
		FullPPROF: conf.Server.EnablePprof,
		rootPath:  "/api/mall",
		conf:      conf,
		checkFunc: checkFunc,
	}
}

func (r *Router) Register(app *gin.Engine) {
	//1.注入鉴权中间件，同时过滤一些不需要鉴权的接口
	app.Use(AuthMiddleware(r.SpanFilter))
	if r.conf.Server.EnablePprof {
		//2.注入pprof
		SetupPprof(app, "/debug/pprof")
	}
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
	root.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hello world",
		})
	})
}
