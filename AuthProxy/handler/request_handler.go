package handler

import (
	"AuthProxyServer/config"
	"AuthProxyServer/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func HandleRequest(c *gin.Context) {

	// 获取环境变量COOKIE_KEYS的值，并按逗号分隔为切片
	cookieKeys := []string{"session"} // 默认只检查"session"键
	if config.EnvCookieKeys != "" {
		cookieKeys = strings.Split(config.EnvCookieKeys, ",")
		for i := range cookieKeys {
			cookieKeys[i] = strings.TrimSpace(cookieKeys[i])
		}
	}
	// 获取代理请求是否为websocket
	connection := c.GetHeader("Connection")

	// 检查 cookie 中的 session 值是否正确
	found := false

	host, tokens, err := utils.GetFrontendHost(c.Request.Host)
	if err != nil {
		log.Printf("Redis error: %s", err)
		if connection == "Upgrade" {
			utils.HandleWebSocketConnectionError(c, "转发服务异常，请联系管理员！")
		} else {
			if config.RedirectURL != "" {
				// 如果存在 REDIRECT_URL 环境变量，则使用配置的重定向链接
				c.Redirect(http.StatusFound, config.RedirectURL)
			} else {
				utils.ServeFileWithStatusCode(c, "./static/redis_error.html", http.StatusServiceUnavailable)
			}
		}
		return
	}

	if host == "" || !utils.IsPortValid(host) {
		log.Printf("host not found, request host: %s，Request from IP: %s", c.Request.Host, utils.GetClientIP(c))
		if connection == "Upgrade" {
			utils.HandleWebSocketConnectionError(c, "无法访问，服务未找到。")
		} else {
			if config.RedirectURL != "" {
				// 如果存在 REDIRECT_URL 环境变量，则使用配置的重定向链接
				c.Redirect(http.StatusNotFound, config.RedirectURL)
			} else {
				utils.ServeFileWithStatusCode(c, "./static/not_found_error.html", http.StatusNotFound)
			}
		}
		return
	}

	if len(tokens) == 0 {
		log.Printf("token not set, request host: %s", c.Request.Host)
		if connection == "Upgrade" {
			utils.HandleWebSocketConnectionError(c, "无法访问，服务未启动或者配置出现错误 . . .")
		} else {
			if config.RedirectURL != "" {
				// 如果存在 REDIRECT_URL 环境变量，则使用配置的重定向链接
				c.Redirect(http.StatusBadRequest, config.RedirectURL)
			} else {
				utils.ServeFileWithStatusCode(c, "./static/connection_error.html", http.StatusBadRequest)
			}
		}
		return
	}

	for _, key := range cookieKeys {
		sessionCookie, err := c.Request.Cookie(key)
		if err == nil && utils.IsinSlice(sessionCookie.Value, tokens) {
			found = true
			break
		}
	}

	if !found {
		if connection == "Upgrade" {
			utils.HandleWebSocketConnectionError(c, "无权限访问！")
		} else {
			if config.RedirectURL != "" {
				// 如果存在 REDIRECT_URL 环境变量，则使用配置的重定向链接
				c.Redirect(http.StatusForbidden, config.RedirectURL)
			} else {
				utils.ServeFileWithStatusCode(c, "./static/authentication_error.html", http.StatusForbidden)
			}
		}
		return
	}

	//host := "192.168.5.249:56912"

	if connection == "Upgrade" {
		//fmt.Println(r.Host, r.URL.Path)
		// Connection 值为 Upgrade，进行ws转发
		HandleWebSocket(c, host)
		return
	}

	// 将请求转发为HTTP请求
	handleHTTP(c, host)
}
