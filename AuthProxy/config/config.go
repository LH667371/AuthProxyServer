package config

import (
	"os"
	"regexp"
)

var (
	EnvCookieKeys = os.Getenv("COOKIE_KEYS")
	RedisAddr     = os.Getenv("REDIS_ADDR")
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	RedisDb       = os.Getenv("REDIS_DB")
	RedirectURL   = os.Getenv("REDIRECT_URL")
	HeartbeatTime = os.Getenv("HEARTBEAT_TIME")
	EnablePprof   = os.Getenv("ENABLE_PPROF")
	UrlPattern    = regexp.MustCompile(`^(https?://)?([^:/]+(:\d+)?)`)
)
