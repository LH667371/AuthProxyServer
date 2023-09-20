package utils

import (
	"AuthProxyServer/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetHeartbeatTime() time.Duration {
	if config.HeartbeatTime == "" {
		return 30 * time.Second
	}
	heartbeatTime, err := strconv.Atoi(config.HeartbeatTime)
	if err != nil {
		// 处理转换错误
		// 这里可以返回默认的 DB 编号或进行其他错误处理
		return 30 * time.Second
	}
	return time.Duration(heartbeatTime) * time.Second
}

// IsURL 判断是否为URL格式
func IsURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

// RemoveHostAndProtocol 删除URL中的主机部分和协议前缀
func RemoveHostAndProtocol(urlStr string, requestHost string) string {
	// 使用正则表达式匹配URL中的主机部分和协议前缀
	urlMatches := config.UrlPattern.FindStringSubmatch(urlStr)
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

func GetClientIP(c *gin.Context) string {
	r := c.Request
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

func HandleWebSocketConnectionError(c *gin.Context, errorMessage string) {
	upgrades := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrades.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// 关闭 conn，以避免资源泄漏
		if conn != nil {
			err := conn.Close()
			if err != nil {
				log.Printf("Failed to upgrade WebSocket connection: %s", err)
				return
			}
		}
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			//logger.Printf("Failed to close WebSocket connection: %s", err)
			return
		}
	}(conn)

	errMsg := []byte(errorMessage)
	err = conn.WriteMessage(websocket.TextMessage, errMsg)
	if err != nil {
		fmt.Println("WebSocket write error:", err)
		return
	}
}

func ServeFileWithStatusCode(c *gin.Context, filePath string, statusCode int) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "File not found",
		})
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close file: %s", err)
		}
	}(file)

	// 设置状态码
	c.Status(statusCode)

	// 将文件内容复制到响应主体
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send file",
		})
		return
	}
}
