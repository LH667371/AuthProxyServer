package handler

import (
	"AuthProxyServer/config"
	"AuthProxyServer/utils"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strings"
)

func handleHTTP(c *gin.Context, host string) {
	defer c.Abort()

	// 构建目标URL，将请求转发到目标主机
	targetURL := fmt.Sprintf("http://%s%s", host, c.Request.URL.Path)

	// 读取原始请求的Body数据
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %s", err)
		// 错误处理
		return
	}

	// 创建新的HTTP请求
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewReader(body))
	if err != nil {
		log.Printf("Failed to create HTTP request: %s", err)
		if config.RedirectURL != "" {
			// 如果存在 REDIRECT_URL 环境变量，则使用配置的重定向链接
			c.Redirect(http.StatusBadRequest, config.RedirectURL)
		} else {
			utils.ServeFileWithStatusCode(c, "./static/connection_error.html", http.StatusBadRequest)
		}
		return
	}

	// 将原始请求的Header复制到新的请求中
	req.URL.RawQuery = c.Request.URL.RawQuery
	// 将原始请求的Header放入新的请求中，去除重复的字段
	req.Header = c.Request.Header
	req.Header.Del("Origin")

	// 发送HTTP请求并获取响应
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 返回一个错误，禁止自动重定向
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send HTTP request: %s", err)

		if config.RedirectURL != "" {
			// 如果存在 REDIRECT_URL 环境变量，则使用配置的重定向链接
			c.Redirect(http.StatusBadRequest, config.RedirectURL)
		} else {
			utils.ServeFileWithStatusCode(c, "./static/connection_error.html", http.StatusBadRequest)
		}
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %s", err)
		}
	}(resp.Body)

	// 将响应的Header复制到原始响应中
	copyHeaders(c.Writer.Header(), resp.Header, host)
	// 将响应的状态码写回给客户端
	c.Status(resp.StatusCode)

	// 将响应的Body复制到原始响应中
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		log.Printf("Failed to copy HTTP response body: %s", err)
		return
	}
}

func copyHeaders(dst http.Header, src http.Header, requestHost string) {
	existingHeaders := make(map[string][]string)

	for key, values := range dst {
		lowerKey := strings.ToLower(key)
		existingHeaders[lowerKey] = values
	}

	for key, values := range src {
		lowerKey := strings.ToLower(key) // 将参数名称转换为小写
		for _, value := range values {
			// 检查字段是否为 Location，并且值是否为 URL 格式
			if lowerKey == "location" && utils.IsURL(value) {
				// 删除 URL 中的主机部分和协议前缀
				value = utils.RemoveHostAndProtocol(value, requestHost)
			}
			// 检查是否存在相同的参数值（不区分大小写）
			if isHeaderExistIgnoreCase(existingHeaders[lowerKey], value) {
				continue
			}

			// 添加响应头参数
			dst.Add(lowerKey, value)
		}
	}
}

// 检查切片中是否存在相同的参数值（不区分大小写）
func isHeaderExistIgnoreCase(slice []string, value string) bool {
	for _, existingValue := range slice {
		if strings.EqualFold(existingValue, value) {
			return true
		}
	}

	return false
}
