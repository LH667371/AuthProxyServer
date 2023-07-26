package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func handleWebSocket(w http.ResponseWriter, r *http.Request, host string) {
	// 构建目标URL，将WebSocket连接转发到目标主机
	targetURL := fmt.Sprintf("ws://%s%s", host, r.URL.Path)

	if r.URL.RawQuery != "" {
		targetURL = fmt.Sprintf("%s?%s", targetURL, r.URL.RawQuery)
	}

	headers := make(http.Header)

	for k, v := range r.Header {
		if k == "Origin" ||
			k == "Upgrade" ||
			k == "Connection" ||
			k == "Sec-Websocket-Key" ||
			k == "Sec-Websocket-Version" ||
			k == "Sec-Websocket-Extensions" {
			//fmt.Println(k, v[0])
		} else {
			headers.Set(k, v[0])
			//fmt.Println("set ==>", k, v[0])
		}
	}

	// 创建到目标服务器的WebSocket连接
	dialer := websocket.DefaultDialer
	targetConn, resp, err := dialer.Dial(targetURL, headers)
	if err != nil {
		logger.Printf("Failed to establish WebSocket connection to target: %s", err)
		return
	}

	if err != nil {
		logger.Printf("set read dead line error: %s", err)
	}

	defer func(targetConn *websocket.Conn) {
		err := targetConn.Close()
		if err != nil {
			logger.Printf("Failed to establish WebSocket connection to target: %s", err)
			return
		}
	}(targetConn)

	//fmt.Println(resp.Header.Get("Sec-Websocket-Accept"))
	respHeaders := make(http.Header)
	for k, v := range resp.Header {
		if k == "Upgrade" ||
			k == "Connection" ||
			k == "Sec-Websocket-Accept" {
			//fmt.Println(k, v[0])
		} else {
			respHeaders.Set(k, v[0])
			//fmt.Println("set ==>", k, v[0])
		}
	}

	conn, err := upgrades.Upgrade(w, r, respHeaders)
	if err != nil {
		logger.Printf("Failed to upgrade WebSocket connection: %s", err)
		return
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Printf("Failed to close WebSocket connection: %s", err)
			return
		}
	}()

	if err != nil {
		logger.Printf("set read dead line error: %s", err)
	}

	//go sendHeartbeat(targetConn)
	go sendHeartbeat(conn)

	// 启动两个 goroutine 来进行双向消息转发
	go copyWebSocketMessages(targetConn, conn)
	copyWebSocketMessages(conn, targetConn)

}

func copyWebSocketMessages(dst *websocket.Conn, src *websocket.Conn) {
	for {
		messageType, message, err := src.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				logger.Printf("Failed to read WebSocket message: %s", err)
			}
			break
		}

		if err := dst.WriteMessage(messageType, message); err != nil {
			logger.Printf("Failed to write WebSocket message: %s", err)
			break
		}
	}
}

func sendHeartbeat(conn *websocket.Conn) {
	heartbeatTime := getHeartbeatTime()

	ticker := time.NewTicker(heartbeatTime) // 每30秒发送一次心跳消息

	for range ticker.C {
		err := conn.WriteMessage(websocket.PingMessage, []byte{})
		if err != nil {
			// 发送心跳消息失败，可能连接已关闭或发生错误
			return
		}
	}
}
