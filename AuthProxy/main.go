package main

import (
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

/*
	COOKIE_KEYS: 鉴权获取的值，Default：session，示例：session1,session2
	REDIS_ADDR: Redis 服务器地址, Default：127.0.0.1:6379
	REDIS_PASSWORD: Redis 访问密码（如果设置了密码），Default：空字符串
	REDIS_DB: Redis 数据库编号，Default：0
	REDIRECT_URL: 开启出现错误重定向的地址， 示例：https://www.xxx.com
	HEARTBEAT_TIME: websocket连接心跳，单位：秒(s)， Default：30s
*/

/*
	// 添加性能分析路由
	go tool pprof http://localhost:6060/debug/pprof/heap
	top：显示占用内存最多的函数。
	list <function>：显示特定函数的源代码和调用位置。
	web：在浏览器中显示可视化的分析结果（需要Graphviz支持）。
*/

func main() {
	// 创建日志文件
	file, err := os.OpenFile("auth-proxy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
	}(file)

	// 设置日志输出到文件和终端
	logger = log.New(io.MultiWriter(file, os.Stdout), "", log.LstdFlags)

	http.HandleFunc("/", handleRequest)
	logger.Fatal(http.ListenAndServe(":80", myMiddleware(http.DefaultServeMux)))
}
