package main

import (
	"AuthProxyServer/config"
	"AuthProxyServer/handler"
	"AuthProxyServer/middlewares"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

/*
	COOKIE_KEYS: 鉴权获取的值，示例：session,training_session
	REDIS_ADDR: Redis 服务器地址, 默认：127.0.0.1:6379
	REDIS_PASSWORD: Redis 访问密码（如果设置了密码），默认：空字符串
	REDIS_DB: Redis 数据库编号，默认：0
	REDIRECT_URL: 开启出现错误重定向的地址， 示例：https://xxxxxx.com
	HEARTBEAT_TIME: websocket连接心跳，单位：秒(s)， 默认 30s
	ENABLE_PPROF: 是否开启pprof，开启设置环境变量为 "true", 重启容器
*/
/*
	_ "net/http/pprof"
	// 添加性能分析路由
	go tool pprof http://localhost/debug/pprof/heap
	top：显示占用内存最多的函数。
	list <function>：显示特定函数的源代码和调用位置。
	web：在浏览器中显示可视化的分析结果（需要Graphviz支持）。
*/

func main() {
	if config.EnablePprof == "true" {
		go func() {
			pprofRouter := gin.Default()
			pprof.Register(pprofRouter)

			err := pprofRouter.Run(":6060")
			if err != nil {
				return
			}
		}()
	}

	router := gin.Default()
	// 日志设置
	config.SetupLogger()
	gin.SetMode(gin.ReleaseMode)

	// 使用CORS中间件
	router.Use(middlewares.CorsMiddleware)
	// 路由注册
	router.Any("/*path", handler.HandleRequest)

	// Default listen and serve on 0.0.0.0:80
	err := router.Run(":80")
	if err != nil {
		return
	}
}
