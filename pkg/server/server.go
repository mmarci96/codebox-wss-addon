package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ClientManager struct {
	clients map[*websocket.Conn]bool
	lock    sync.Mutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (cm *ClientManager) Add(conn *websocket.Conn) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.clients[conn] = true
	log.Printf("Client connected. Total: %d", len(cm.clients))
}

func (cm *ClientManager) Remove(conn *websocket.Conn) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	delete(cm.clients, conn)
	log.Printf("Client disconnected. Total: %d", len(cm.clients))
}

func (cm *ClientManager) Broadcast(message string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	for conn := range cm.clients {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println("Write error:", err)
			conn.Close()
			delete(cm.clients, conn)
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins for dev/test
	},
}

func wsHandler(cm *ClientManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return
		}
		cm.Add(conn)
		defer cm.Remove(conn)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received: %s", msg)
			cm.Broadcast(fmt.Sprintf("Echo: %s", msg))
		}
	}
}

func Run() *gin.Engine {
	r := gin.Default()
	cm := NewClientManager()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})
	r.GET("/ws", wsHandler(cm))

	return r
}
