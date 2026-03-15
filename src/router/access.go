package router

import (
	"bytes"
	"io"
	"time"

	"github.com/Sam-Stranding/SamMall/src/consts"
	"github.com/Sam-Stranding/SamMall/src/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetRequestBody(c *gin.Context) string {
	data, _ := io.ReadAll(c.Request.Body)
	return string(data)
}

func GetResponseBody(c *gin.Context) string {
	resp := c.Request.Response
	if resp == nil || resp.Body == nil {
		return ""
	}
	data, _ := io.ReadAll(c.Request.Response.Body)
	return string(data)
}

type responseWriterWrapper struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func AccessLogMiddleware(filter func(*gin.Context) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if filter != nil && !filter(c) {
			c.Next()
			return
		}
		body := GetRequestBody(c)
		c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(body)))
		begin := time.Now()
		fields := []zap.Field{
			zap.String("ip", c.RemoteIP()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("body", body),
			zap.String("token", c.GetHeader(consts.UserTokenKey)),
		}

		var responseBody bytes.Buffer
		multiWriter := io.MultiWriter(c.Writer, &responseBody)
		c.Writer = &responseWriterWrapper{
			ResponseWriter: c.Writer,
			Writer:         multiWriter,
		}

		//日志中间件,执行剩余中间件链的方法
		c.Next()
		respBody := responseBody.String()
		if len(respBody) > 1024 {
			respBody = respBody[:1024]
		}
		//zap.Duration是go语言zap日志库的一个字段构造函数，用于记录时间间隔类型的日志字段
		fields = append(fields, zap.Int64("dur_ms", time.Since(begin).Milliseconds()))
		fields = append(fields, zap.Int("status", c.Writer.Status()))
		fields = append(fields, zap.String("resp", respBody))
		logger.Info("access_log", fields...)
	}
}
