package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	Type    int
	Payload []byte
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, host string, wsWg *sync.WaitGroup) {
	defer wsWg.Done()

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
		// 关闭 targetConn，以避免资源泄漏
		if targetConn != nil {
			err := targetConn.Close()
			if err != nil {
				logger.Printf("Failed to establish WebSocket connection to target: %s", err)
				return
			}
		}
		return
	}

	defer func(targetConn *websocket.Conn) {
		err := targetConn.Close()
		if err != nil {
			//logger.Printf("Failed to establish WebSocket connection to target: %s", err)
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

	upgrades := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrades.Upgrade(w, r, respHeaders)
	if err != nil {
		// 关闭 conn，以避免资源泄漏
		if conn != nil {
			err := conn.Close()
			if err != nil {
				logger.Printf("Failed to upgrade WebSocket connection: %s", err)
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

	wg := &sync.WaitGroup{}
	wg.Add(2) // 如果还有其他的协程，可以继续添加

	targetConnMessages := make(chan Message)
	connMessages := make(chan Message)

	//go sendHeartbeat(targetConnMessages)
	go sendHeartbeat(connMessages)

	// 启动两个 goroutine 来进行双向消息转发
	go copyWebSocketMessages(targetConn, conn, targetConnMessages, wg)
	go copyWebSocketMessages(conn, targetConn, connMessages, wg)

	// 等待所有协程完成
	wg.Wait()
}

func copyWebSocketMessages(dst *websocket.Conn, src *websocket.Conn, messages chan Message, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		defer close(messages)

		for {
			messageType, message, err := src.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(
					err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived,
					websocket.CloseAbnormalClosure,
				) {
					logger.Printf("Failed to read WebSocket message: %s", err)
				}
				return
			}

			messages <- Message{Type: messageType, Payload: message}
		}
	}()

	// 写入消息
	for msg := range messages {
		if err := dst.WriteMessage(msg.Type, msg.Payload); err != nil {
			if !errors.Is(err, websocket.ErrCloseSent) {
				logger.Printf("Failed to write WebSocket message: %s", err)
			}
			return
		}
	}
}

func sendHeartbeat(messages chan Message) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	heartbeatTime := getHeartbeatTime()
	ticker := time.NewTicker(heartbeatTime)
	defer ticker.Stop()

	for range ticker.C {
		messages <- Message{Type: websocket.PingMessage, Payload: []byte{}}
	}
}
