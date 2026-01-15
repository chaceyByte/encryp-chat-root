package websocket

import (
	"encoding/json"
	"log"
	"time"

	"chat-server/pkg/types"

	"github.com/gorilla/websocket"
)

// WSConnection WebSocket 连接
type WSConnection struct {
	Conn *websocket.Conn
	IP   string
	hub  *Hub
	Send chan []byte

	// 配置
	readTimeout    time.Duration
	writeTimeout   time.Duration
	pingPeriod     time.Duration
	maxMessageSize int64
}

// NewWSConnection 创建 WebSocket 连接
func NewWSConnection(conn *websocket.Conn, ip string, cfg *types.ChatServerConfig) *WSConnection {
	return &WSConnection{
		Conn:           conn,
		IP:             ip,
		Send:           make(chan []byte, 256),
		readTimeout:    cfg.ReadTimeout,
		writeTimeout:   cfg.WriteTimeout,
		pingPeriod:     cfg.PingPeriod,
		maxMessageSize: cfg.MaxMessageSize,
	}
}

// ReadPump 读取消息泵
func (ws *WSConnection) ReadPump(hub *Hub) {
	defer func() {
		ws.Conn.Close()
	}()

	// 设置读取限制
	ws.Conn.SetReadLimit(ws.maxMessageSize)

	// 设置读取超时
	ws.Conn.SetReadDeadline(time.Now().Add(ws.readTimeout))

	// 设置 Pong 处理
	ws.Conn.SetPongHandler(func(string) error {
		ws.Conn.SetReadDeadline(time.Now().Add(ws.readTimeout))
		return nil
	})

	// 循环读取消息
	for {
		_, messageData, err := ws.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket 读取错误: %v", err)
			}
			break
		}

		// 解析消息
		msg, err := parseMessage(messageData)
		if err != nil {
			log.Printf("消息解析错误: %v", err)
			ws.SendError("消息格式错误")
			continue
		}

		// 处理消息
		ws.handleMessage(hub, msg)
	}
}

// WritePump 写入消息泵
func (ws *WSConnection) WritePump() {
	ticker := time.NewTicker(ws.pingPeriod)
	defer func() {
		ticker.Stop()
		ws.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-ws.Send:
			// 设置写入超时
			ws.Conn.SetWriteDeadline(time.Now().Add(ws.writeTimeout))
			if !ok {
				// 通道已关闭
				ws.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 写入消息
			w, err := ws.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("创建写入器失败: %v", err)
				return
			}
			w.Write(message)

			// 排队其他消息
			n := len(ws.Send)
			for i := 0; i < n; i++ {
				w.Write(<-ws.Send)
			}

			if err := w.Close(); err != nil {
				log.Printf("关闭写入器失败: %v", err)
				return
			}

		case <-ticker.C:
			// 发送 Ping
			ws.Conn.SetWriteDeadline(time.Now().Add(ws.writeTimeout))
			if err := ws.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("发送 Ping 失败: %v", err)
				return
			}
		}
	}
}

// handleMessage 处理消息
func (ws *WSConnection) handleMessage(hub *Hub, msg *types.WebSocketMessage) {
	switch msg.Type {
	case "join":
		if msg.UserID == "" || msg.Username == "" || msg.RoomKey == "" {
			ws.SendError("加入房间消息缺少必要字段")
			return
		}
		hub.HandleJoin(ws, msg.UserID, msg.Username, msg.RoomKey)

	case "leave":
		if msg.UserID == "" || msg.RoomKey == "" {
			return
		}
		hub.HandleLeave(msg.UserID, msg.RoomKey)

	case "message":
		if msg.UserID == "" || msg.Content == "" || msg.RoomKey == "" {
			ws.SendError("聊天消息缺少必要字段")
			return
		}
		hub.HandleMessage(msg.UserID, msg.RoomKey, msg.Content, msg.MessageType, msg.Encrypted)

	default:
		ws.SendError("未知的消息类型")
	}
}

// SendMessage 发送消息
func (ws *WSConnection) SendMessage(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case ws.Send <- data:
		return nil
	default:
		return err
	}
}

// SendError 发送错误消息
func (ws *WSConnection) SendError(errMsg string) {
	msg := map[string]interface{}{
		"type":  "error",
		"error": errMsg,
	}
	data, _ := json.Marshal(msg)
	select {
	case ws.Send <- data:
	default:
	}
}

// parseMessage 解析消息
func parseMessage(data []byte) (*types.WebSocketMessage, error) {
	var msg types.WebSocketMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
