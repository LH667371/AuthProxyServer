package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net"
	"net/url"
	"strconv"
)

func getRedisAddr() string {
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	return redisAddr
}

func getRedisDB() int {
	if redisDb == "" {
		return 0
	}
	db, err := strconv.Atoi(redisDb)
	if err != nil {
		// 处理转换错误
		// 这里可以返回默认的 DB 编号或进行其他错误处理
		return 0
	}
	return db
}

func GetFrontendHost(host string) (string, []string, error) {
	var fromURL string
	var tokens []string

	// 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     getRedisAddr(), // Redis 服务器地址
		Password: redisPassword,  // Redis 访问密码（如果设置了密码）
		DB:       getRedisDB(),   // Redis 数据库编号
	})

	// 检查 Redis 连接是否成功
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return fromURL, tokens, fmt.Errorf("failed to connect to Redis: %s", err)
	}

	// 从 Redis 中获取 frontend 列表
	listName := fmt.Sprintf("frontend:%s", host)
	frontendList, err := client.LRange(context.Background(), listName, 0, -1).Result()
	if err != nil {
		return fromURL, tokens, fmt.Errorf("failed to retrieve frontend list from Redis: %s", err)
	}

	for i := range frontendList {
		fromURL, err = extractHostFromURL(frontendList[i])
		if err != nil {
			return fromURL, tokens, err
		}
	}

	// 获取Set的值
	// 从 Redis 中获取 frontend 列表
	setName := fmt.Sprintf("token:%s", host)
	tokens, err = client.SMembers(context.Background(), setName).Result()
	if err != nil {
		// 处理错误
		return fromURL, tokens, fmt.Errorf("failed to get values from Redis set: %s", err)
	}

	// 关闭Redis连接
	err = client.Close()
	if err != nil {
		// 处理错误
		return fromURL, tokens, fmt.Errorf("failed to close Redis connection: %s", err)
	}

	return fromURL, tokens, nil
}

func extractHostFromURL(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}

func IsPortValid(host string) bool {
	_, strPort, _ := net.SplitHostPort(host)
	// 将端口号字符串转换为整数
	port, err := strconv.Atoi(strPort)
	if err != nil {
		fmt.Println("Failed to convert port to integer:", err)
		return false
	}
	return port >= 1 && port <= 65535
}

func isInSlice(s string, slice []string) bool {
	for _, item := range slice {
		if item == s || item == "ALL" {
			return true
		}
	}
	return false
}
