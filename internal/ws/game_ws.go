package ws

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zhi-minn/dota-sketch-backend/internal/models"
	"github.com/zhi-minn/dota-sketch-backend/internal/services"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type BroadcastMsg struct {
	lobbyId string
	data    []byte
}

type Hub struct {
	lobbies    map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan BroadcastMsg
	mu         sync.RWMutex
}

type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	lobbyId string
	send    chan []byte
}

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WsHandler struct {
	gameService *services.GameService
}

func NewWsHandler(gameService *services.GameService) *WsHandler {
	return &WsHandler{
		gameService: gameService,
	}
}

func NewHub() *Hub {
	return &Hub{
		lobbies:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan BroadcastMsg),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			lobbyId := client.lobbyId
			if _, ok := h.lobbies[lobbyId]; !ok {
				h.lobbies[lobbyId] = make(map[*Client]bool)
			}
			h.lobbies[lobbyId][client] = true
		case client := <-h.unregister:
			lobbyId := client.lobbyId
			if clients, ok := h.lobbies[lobbyId]; ok {
				delete(h.lobbies, lobbyId)
				close(client.send)
				if len(clients) == 0 {
					delete(h.lobbies, lobbyId)
				}
			}
		case lobbyMsg := <-h.broadcast:
			lobbyId := lobbyMsg.lobbyId
			if clients, ok := h.lobbies[lobbyId]; ok {
				for client := range clients {
					select {
					case client.send <- lobbyMsg.data:
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}
		}
	}
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("readPump err:", err)
			break
		}

		log.Println(message)
	}
}

func (c *Client) writePump(h *Hub) {
	defer c.conn.Close()

	for {
		msg, ok := <-c.send
		if !ok {
			log.Println("Send channel closed, disconnecting client")
			return
		}

		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("writePump err:", err)
			return
		}
	}
}

func processMessage(client *Client, rawMsg []byte) {
	var msg Message
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		log.Println("Invalid msg:", string(rawMsg))
		return
	}

	switch msg.Type {
	case "DRAW_EVENT":
	case "WORD_GUESS":
	case "LOBBY_UPDATE":
	default:
		log.Println("Unknown message type:", msg.Type)
	}
}

func (gs *WsHandler) ServeWs(h *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		lobbyId := c.Query("lobbyCode")
		client := &Client{
			hub:     h,
			conn:    conn,
			lobbyId: lobbyId,
			send:    make(chan []byte, 256),
		}

		h.register <- client

		go sendLobbySettings(client)

		go client.readPump(h)
		go client.writePump(h)
	}
}

func sendLobbySettings(c *Client) {
	lobbyId := c.lobbyId

	models.Mutex.Lock()
	game, exists := models.ActiveGames[lobbyId]
	models.Mutex.Unlock()

	if !exists {
		log.Println("Lobby settings not found", lobbyId)
		return
	}

	settingsJSON, err := json.Marshal(game.Settings)
	if err != nil {
		log.Println("Error marshalling settings:", err)
		return
	}

	message := Message{
		Type:    "LOBBY_SETTINGS",
		Payload: settingsJSON,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}

	c.send <- messageJSON
}
