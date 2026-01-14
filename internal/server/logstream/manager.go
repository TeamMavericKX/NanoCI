package logstream

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type LogManager struct {
	rdb         *redis.Client
	subscribers map[string]map[*websocket.Conn]bool
	mu          sync.RWMutex
}

func NewLogManager(rdb *redis.Client) *LogManager {
	return &LogManager{
		rdb:         rdb,
		subscribers: make(map[string]map[*websocket.Conn]bool),
	}
}

func (m *LogManager) HandleWS(w http.ResponseWriter, r *http.Request, buildID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.L().Error("ws upgrade failed", zap.Error(err))
		return
	}
	defer conn.Close()

	m.mu.Lock()
	if m.subscribers[buildID] == nil {
		m.subscribers[buildID] = make(map[*websocket.Conn]bool)
		// Start a redis subscriber for this build if it's the first client
		go m.listenRedis(buildID)
	}
	m.subscribers[buildID][conn] = true
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		delete(m.subscribers[buildID], conn)
		if len(m.subscribers[buildID]) == 0 {
			delete(m.subscribers, buildID)
		}
		m.mu.Unlock()
	}()

	// Keep connection alive/wait for close
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (m *LogManager) listenRedis(buildID string) {
	pubsub := m.rdb.Subscribe(context.Background(), fmt.Sprintf("logs:%s", buildID))
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		m.mu.RLock()
		conns := m.subscribers[buildID]
		for conn := range conns {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				zap.L().Error("failed to write to ws", zap.Error(err))
			}
		}
		m.mu.RUnlock()

		// Stop listening if no more subscribers
		m.mu.RLock()
		_, exists := m.subscribers[buildID]
		m.mu.RUnlock()
		if !exists {
			return
		}
	}
}
