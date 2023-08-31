package main

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getHeartbeatTime() time.Duration {
	if HeartbeatTime == "" {
		return 30 * time.Second
	}
	heartbeatTime, err := strconv.Atoi(HeartbeatTime)
	if err != nil {
		// 处理转换错误
		// 这里可以返回默认的 DB 编号或进行其他错误处理
		return 30 * time.Second
	}
	return time.Duration(heartbeatTime) * time.Second
}

// 判断是否为URL格式
func isURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

// 删除URL中的主机部分和协议前缀
func removeHostAndProtocol(urlStr string, requestHost string) string {
	// 使用正则表达式匹配URL中的主机部分和协议前缀
	re := regexp.MustCompile(`^(https?://)?([^:/]+(:\d+)?)`)
	urlMatches := re.FindStringSubmatch(urlStr)
	if len(urlMatches) >= 3 {
		//protocol := urlMatches[1] // 匹配到的协议部分（包括 "http://" 或 "https://"）
		host := urlMatches[2] // 匹配到的主机部分

		// 如果主机部分和请求头的Host一致，则删除主机部分和协议前缀
		if host == requestHost {
			urlStr = strings.Replace(urlStr, urlMatches[0], "", 1)
		}
	}

	return urlStr
}

func getClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ", ")
		// 第一个IP地址是客户端的真实IP
		return ips[0]
	}
	// 如果 X-Forwarded-For 不存在，则使用 RemoteAddr
	remoteAddrParts := strings.Split(r.RemoteAddr, ":")
	if len(remoteAddrParts) > 0 {
		return remoteAddrParts[0]
	}
	// 如果都无法获取，则返回空字符串
	return ""
}
