package main

import (
	"net/http"
	"strings"
	"sync"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 设置允许所有请求跨域
	if r.Header.Get("Origin") != "" {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Disposition")
	}

	// 获取环境变量COOKIE_KEYS的值，并按逗号分隔为切片
	cookieKeys := []string{"session"} // 默认只检查"session"键
	if envCookieKeys != "" {
		cookieKeys = strings.Split(envCookieKeys, ",")
		for i := range cookieKeys {
			cookieKeys[i] = strings.TrimSpace(cookieKeys[i])
		}
	}

	// 检查 cookie 中的 session 值是否正确
	found := false

	host, tokens, err := GetFrontendHost(r.Host)
	if err != nil {
		logger.Printf("Redis error: %s", err)
		if redirectURL != "" {
			// 如果存在 REDIRECT_URL 环境变量，则使用配置的重定向链接
			http.Redirect(w, r, redirectURL, http.StatusFound)
		} else {
			http.ServeFile(w, r, "./static/redis_error.html")
		}
		return
	}

	if host == "" || !IsPortValid(host) {
		logger.Printf("host not found, request host: %s，Request from IP: %s", r.Host, getClientIP(r))
		if redirectURL != "" {
			// 如果存在 REDIRECT_URL 环境变量，则使用配置的重定向链接
			http.Redirect(w, r, redirectURL, http.StatusFound)
		} else {
			http.ServeFile(w, r, "./static/not_found_error.html")
		}
		return
	}

	if len(tokens) == 0 {
		logger.Printf("token not set, request host: %s", r.Host)
		if redirectURL != "" {
			// 如果存在 REDIRECT_URL 环境变量，则使用配置的重定向链接
			http.Redirect(w, r, redirectURL, http.StatusFound)
		} else {
			http.ServeFile(w, r, "./static/connection_error.html")
		}
		return
	}

	for _, key := range cookieKeys {
		sessionCookie, err := r.Cookie(key)
		if err == nil && isInSlice(sessionCookie.Value, tokens) {
			found = true
			break
		}
	}

	if r.Method == http.MethodOptions {
		found = true
	}

	if !found {
		if redirectURL != "" {
			// 如果存在 REDIRECT_URL 环境变量，则使用配置的重定向链接
			http.Redirect(w, r, redirectURL, http.StatusFound)
		} else {
			http.ServeFile(w, r, "./static/authentication_error.html")
		}
		return
	}

	//host := "192.168.5.249:56912"

	wg := &sync.WaitGroup{}

	connection := r.Header.Get("Connection")
	if connection == "Upgrade" {
		//fmt.Println(r.Host, r.URL.Path)
		// Connection 值为 Upgrade，进行ws转发
		wg.Add(1)
		go handleWebSocket(w, r, host, wg)
		wg.Wait()
		return
	}

	// 将请求转发为HTTP请求
	wg.Add(1)
	go handleHTTP(w, r, host, wg)
	wg.Wait()
}
