package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var (
	envCookieKeys = os.Getenv("COOKIE_KEYS")
	redisAddr     = os.Getenv("REDIS_ADDR")
	redisPassword = os.Getenv("REDIS_PASSWORD")
	redisDb       = os.Getenv("REDIS_DB")
	redirectURL   = os.Getenv("REDIRECT_URL")
	HeartbeatTime = os.Getenv("HEARTBEAT_TIME")
)

var (
	logger *log.Logger
)

var upgrades = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
