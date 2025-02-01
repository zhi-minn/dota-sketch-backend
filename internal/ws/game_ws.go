package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	lobbyId string
	send    chan []byte
}

type Hub struct {
	lobbies    map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan BroadcastMsg
	mu         sync.RWMutex
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

func ServeWs(h *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		client := &Client{
			hub:  h,
			conn: conn,
			send: make(chan []byte, 256),
		}

		h.register <- client
	}
}
