package config

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func SetupLogger() {
	// 打开或创建日志文件
	logFilePath := "gin_app.log"
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件：%v", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	// 设置日志输出
	gin.DefaultWriter = logFile
	gin.DefaultErrorWriter = logFile
}
