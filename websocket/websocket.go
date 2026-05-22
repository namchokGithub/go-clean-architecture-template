package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client holds a WebSocket connection for a single user.
type Client struct {
	UserID string
	Conn   *websocket.Conn
}

// Manager tracks connected clients and named groups.
type Manager struct {
	mu      sync.RWMutex
	clients map[string]*Client
	groups  map[string]map[string]bool
}

// NewManager creates an empty WebSocket manager.
func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]*Client),
		groups:  make(map[string]map[string]bool),
	}
}

// Register adds a client to the manager.
func (m *Manager) Register(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[c.UserID] = c
}

// Unregister removes a client and closes the connection.
func (m *Manager) Unregister(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c, ok := m.clients[userID]; ok {
		c.Conn.Close()
		delete(m.clients, userID)
	}
}

// JoinGroup adds a user to a named group.
func (m *Manager) JoinGroup(groupName, userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.groups[groupName] == nil {
		m.groups[groupName] = make(map[string]bool)
	}
	m.groups[groupName][userID] = true
}

// SendToUsers broadcasts data to a list of users.
func (m *Manager) SendToUsers(userIDs []string, data []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, id := range userIDs {
		if c, ok := m.clients[id]; ok {
			_ = c.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

// RefreshProduct broadcasts a product refresh event to the "product" group.
func (m *Manager) RefreshProduct(productID int64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	group, ok := m.groups["product"]
	if !ok {
		return
	}
	msg := []byte(`{"event":"product_refresh"}`)
	for userID := range group {
		if c, ok2 := m.clients[userID]; ok2 {
			_ = c.Conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
	_ = productID
}
