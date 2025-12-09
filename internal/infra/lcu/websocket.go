package lcu

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
)

// EventType WebSocket 事件类型
type EventType string

const (
	EventCreate EventType = "Create"
	EventUpdate EventType = "Update"
	EventDelete EventType = "Delete"
)

// LCUEvent LCU WebSocket 事件
type LCUEvent struct {
	URI       string          `json:"uri"`
	EventType EventType       `json:"eventType"`
	Data      json.RawMessage `json:"data"`
}

// EventHandler 事件处理器函数类型
type EventHandler func(event *LCUEvent) error

// EventSubscription 事件订阅
type EventSubscription struct {
	URI     string       // 订阅的URI模式
	Types   []EventType  // 订阅的事件类型
	Handler EventHandler // 事件处理器
}

// WebSocketManager WebSocket 管理器
type WebSocketManager struct {
	client        *Client
	conn          *websocket.Conn
	subscriptions []EventSubscription
	mu            sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	log           *log.Helper
	reconnecting  bool
}

// NewWebSocketManager 创建 WebSocket 管理器
func NewWebSocketManager(client *Client, logger log.Logger) *WebSocketManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketManager{
		client:        client,
		subscriptions: make([]EventSubscription, 0),
		ctx:           ctx,
		cancel:        cancel,
		log:           log.NewHelper(log.With(logger, "module", "lcu/websocket")),
	}
}

// Subscribe 订阅事件
func (wsm *WebSocketManager) Subscribe(uri string, types []EventType, handler EventHandler) {
	wsm.mu.Lock()
	defer wsm.mu.Unlock()

	wsm.subscriptions = append(wsm.subscriptions, EventSubscription{
		URI:     uri,
		Types:   types,
		Handler: handler,
	})

	wsm.log.Infof("Subscribed to URI: %s, types: %v", uri, types)
}

// Start 启动 WebSocket 连接
func (wsm *WebSocketManager) Start() error {
	wsm.client.Mutex().RLock()
	port := wsm.client.Port
	token := wsm.client.Token
	connected := wsm.client.Connected
	wsm.client.Mutex().RUnlock()

	if !connected {
		return fmt.Errorf("LCU client is not connected")
	}

	url := fmt.Sprintf("wss://127.0.0.1:%d/", port)

	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	header := http.Header{}
	header.Add("Authorization", "Basic "+basicAuth("riot", token))

	conn, _, err := dialer.Dial(url, header)
	if err != nil {
		return fmt.Errorf("failed to connect to LCU WebSocket: %w", err)
	}

	wsm.conn = conn
	wsm.log.Info("WebSocket connected to LCU")

	// 订阅所有事件
	if err := wsm.subscribeEvents(); err != nil {
		conn.Close()
		return fmt.Errorf("failed to subscribe events: %w", err)
	}

	// 启动消息处理协程
	go wsm.handleMessages()

	return nil
}

// subscribeEvents 订阅所有已注册的事件
func (wsm *WebSocketManager) subscribeEvents() error {
	wsm.mu.RLock()
	defer wsm.mu.RUnlock()

	// 订阅 OnJsonApiEvent 以接收所有事件
	msg := []interface{}{5, "OnJsonApiEvent"}
	if err := wsm.conn.WriteJSON(msg); err != nil {
		return err
	}

	wsm.log.Info("Subscribed to OnJsonApiEvent")
	return nil
}

// handleMessages 处理接收到的消息
func (wsm *WebSocketManager) handleMessages() {
	defer func() {
		if wsm.conn != nil {
			wsm.conn.Close()
		}
	}()

	for {
		select {
		case <-wsm.ctx.Done():
			wsm.log.Info("WebSocket message handler stopped")
			return
		default:
			_, message, err := wsm.conn.ReadMessage()
			if err != nil {
				wsm.log.Errorf("Failed to read message: %v", err)

				// 尝试重连
				if !wsm.reconnecting {
					go wsm.reconnect()
				}
				return
			}

			// 解析消息
			var rawMsg []json.RawMessage
			if err := json.Unmarshal(message, &rawMsg); err != nil {
				wsm.log.Errorf("Failed to unmarshal message: %v", err)
				continue
			}

			// 消息格式: [opcode, event, data]
			if len(rawMsg) < 3 {
				continue
			}

			var event LCUEvent
			if err := json.Unmarshal(rawMsg[2], &event); err != nil {
				wsm.log.Errorf("Failed to unmarshal event: %v", err)
				continue
			}

			// 分发事件
			go wsm.dispatchEvent(&event)
		}
	}
}

// dispatchEvent 分发事件到订阅者
func (wsm *WebSocketManager) dispatchEvent(event *LCUEvent) {
	wsm.mu.RLock()
	defer wsm.mu.RUnlock()

	for _, sub := range wsm.subscriptions {
		// 检查 URI 是否匹配
		if sub.URI == "" || matchURI(sub.URI, event.URI) {
			// 检查事件类型是否匹配
			if len(sub.Types) == 0 || containsEventType(sub.Types, event.EventType) {
				// 调用处理器
				if err := sub.Handler(event); err != nil {
					wsm.log.Errorf("Event handler error for URI %s: %v", event.URI, err)
				}
			}
		}
	}
}

// reconnect 重新连接 WebSocket
func (wsm *WebSocketManager) reconnect() {
	wsm.reconnecting = true
	defer func() {
		wsm.reconnecting = false
	}()

	wsm.log.Info("Attempting to reconnect WebSocket...")

	for i := 0; i < 10; i++ {
		select {
		case <-wsm.ctx.Done():
			return
		default:
			time.Sleep(time.Second * 3)

			if err := wsm.Start(); err != nil {
				wsm.log.Errorf("Reconnect attempt %d failed: %v", i+1, err)
				continue
			}

			wsm.log.Info("WebSocket reconnected successfully")
			return
		}
	}

	wsm.log.Error("Failed to reconnect WebSocket after 10 attempts")
}

// Stop 停止 WebSocket 连接
func (wsm *WebSocketManager) Stop() {
	wsm.cancel()
	if wsm.conn != nil {
		wsm.conn.Close()
	}
	wsm.log.Info("WebSocket stopped")
}

// matchURI 检查URI是否匹配
func matchURI(pattern, uri string) bool {
	// 简单的前缀匹配
	// 可以扩展为更复杂的模式匹配
	if pattern == uri {
		return true
	}

	// 支持通配符
	if pattern == "*" {
		return true
	}

	return false
}

// containsEventType 检查事件类型数组是否包含指定类型
func containsEventType(types []EventType, eventType EventType) bool {
	for _, t := range types {
		if t == eventType {
			return true
		}
	}
	return false
}

// basicAuth 生成 Basic Auth 字符串
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
