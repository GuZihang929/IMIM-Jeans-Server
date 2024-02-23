package system

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type WebSocketController struct {
}

func (ws WebSocketController) Ws(c *gin.Context) {
	fmt.Println("123")
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// 允许所有来源的 WebSocket 连接
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	// 处理 WebSocket 连接
	run(conn)
}

func run(c *websocket.Conn) {
	go readMessage(c)
	go writeMessage(c)
}

// readMessage 开始从 Connection 中读取消息
func readMessage(c *websocket.Conn) {

	ticker := time.NewTicker(3 * time.Second) // 创建一个每3秒触发一次的定时器

	// 定时器会定期触发，直到程序退出或定时器停止
	for {
		select {
		case <-ticker.C: // 定时器触发时执行的操作
			messageType, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Error during message reading:", err)
				break
			}
			log.Printf("Received messageType: %d", messageType)
			log.Printf("Received: %s", message)
		}
	}
}

// writeMessage 开始向 Connection 中写入消息队列中的消息
func writeMessage(c *websocket.Conn) {
	str := "Hello, World!"
	bytes := []byte(str)
	ticker := time.NewTicker(3 * time.Second) // 创建一个每3秒触发一次的定时器

	// 定时器会定期触发，直到程序退出或定时器停止
	for {
		select {
		case <-ticker.C: // 定时器触发时执行的操作
			err := c.WriteMessage(1, bytes)
			if err != nil {
				log.Println("Error during message writing:", err)
				break
			}
		}
	}

}
