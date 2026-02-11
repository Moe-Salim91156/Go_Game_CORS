package handlers

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Hub struct {
	Rooms map[string]map[*websocket.Conn]bool
	// pretty good way but will take me a while to grasp it visually in my brain,
	// but i can use it
	mu sync.Mutex
	// ofcourse , mutexes to prevent two stupid players of modifying the Hub as the same time

}

var GameHub = &Hub{
	Rooms: make(map[string]map[*websocket.Conn]bool),
}

func (h *Hub) RegisterIntoRooms(roomID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.Rooms[roomID] == nil {
		h.Rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	h.Rooms[roomID][conn] = true

}

func (h *Hub) Unregister(roomID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if connection, exists := h.Rooms[roomID]; exists {
		delete(connection, conn)
		if len(connection) == 0 {
			delete(h.Rooms, roomID)
		}
	}
}

// THE SHOUTING , update on every side
func (h *Hub) BroadCast(roomID string, data interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.Rooms[roomID] {
		err := conn.WriteJSON(data)
		if err != nil {
			defer conn.Close()
			delete(h.Rooms[roomID], conn)
		}
	}
}
