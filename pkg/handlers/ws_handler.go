package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	username string
	isRoot   bool
}

type ClientManager struct {
	clients []*Client
	lock    sync.Mutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: []*Client{},
	}
}

func (cm *ClientManager) Add(client *Client) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.clients = append(cm.clients, client)
	log.Printf("Client connected. Total: %d", len(cm.clients))
}

func (cm *ClientManager) Remove(conn *websocket.Conn) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	for i, c := range cm.clients {
		if c.conn == conn {
			c.conn.Close()
			cm.clients = append(cm.clients[:i], cm.clients[i+1:]...)
			break
		}
	}
	log.Printf("Client disconnected. Total: %d", len(cm.clients))
}

func (cm *ClientManager) FindRoot() *Client {
	for _, c := range cm.clients {
		if c.isRoot {
			return c
		}
	}
	return nil
}

func (cm *ClientManager) FindByUsername(username string) *Client {
	for _, c := range cm.clients {
		if c.username == username {
			return c
		}
	}
	return nil
}

func (cm *ClientManager) Broadcast(message string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	aliveClients := []*Client{}

	for _, client := range cm.clients {
		err := client.conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Broadcast write error:", err)
			client.conn.Close()
			continue
		}
		aliveClients = append(aliveClients, client)
	}

	cm.clients = aliveClients
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(cm *ClientManager, hostUser string) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade failed:", err)
			return
		}
		defer conn.Close()

		client := &Client{
			conn: conn,
		}
		cm.Add(client)
		defer cm.Remove(conn)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}

			var parsed map[string]any
			if err := json.Unmarshal(msg, &parsed); err != nil {
				log.Println("Invalid JSON:", string(msg))
				continue
			}

			switch parsed["type"] {
			case "login":
				username := parsed["username"].(string)
				isRoot := username == hostUser
				client.username = username
				client.isRoot = isRoot

				response := map[string]any{
					"type":     "login_ack",
					"username": username,
					"root":     isRoot,
				}
				respJSON, _ := json.Marshal(response)
				conn.WriteMessage(websocket.TextMessage, respJSON)

			case "getsecret":
				if client.isRoot {
					continue
				}
				root := cm.FindRoot()
				if root != nil {
					request := map[string]any{
						"type":     "approve_secret",
						"username": client.username,
					}
					reqJSON, _ := json.Marshal(request)
					root.conn.WriteMessage(websocket.TextMessage, reqJSON)
				} else {
					errResp := map[string]any{
						"type":  "error",
						"error": "No root user connected",
					}
					errJSON, _ := json.Marshal(errResp)
					conn.WriteMessage(websocket.TextMessage, errJSON)
				}

			case "secret_response":
				targetUser := parsed["username"].(string)
				secret := parsed["secret"].(string)
				targetClient := cm.FindByUsername(targetUser)
				if targetClient != nil {
					result := map[string]any{
						"type":   "secret_result",
						"secret": secret,
					}
					resultJSON, _ := json.Marshal(result)
					targetClient.conn.WriteMessage(websocket.TextMessage, resultJSON)
				}

			default:
				log.Println("Unhandled type:", parsed["type"])
			}
		}
	}
}
