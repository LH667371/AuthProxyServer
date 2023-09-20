package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CorsMiddleware(c *gin.Context) {
	// 设置CORS头部
	c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Disposition")

	// 如果是预检请求（OPTIONS请求），则返回空响应
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	// 继续处理其他请求
	c.Next()
}
