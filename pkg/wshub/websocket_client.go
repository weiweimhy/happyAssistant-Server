package wshub

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type int
	Data []byte
}

// IClient WebSocket客户端接口
type IClient interface {
	SendText(msg []byte) error
	SendBinary(msg []byte) error
	Close()
	GetBaseClient() *Client
	// Context相关方法
	SetContextValue(key string, value interface{})
	GetContextValue(key string) interface{}
	GetContextString(key string) string
	GetContextInt(key string) int
	GetContextBool(key string) bool
}

type Client struct {
	OnMessage  func(client IClient, msg []byte)
	OnError    func(client IClient, err error)
	OnClose    func(client IClient)
	conn       *websocket.Conn
	sendBuffer chan *Message
	closeOnce  sync.Once
	ctx        context.Context
	ctxMu      sync.Mutex
	*ClientConfig
}

// 确保Client实现IClient接口
var _ IClient = (*Client)(nil)

type ClientConfig struct {
	chanLength    int
	readDeadline  time.Duration
	writeDeadline time.Duration
	readLimit     int64
	supportPing   bool
	pingPeriod    time.Duration
}

type ClientOption func(*ClientConfig)

func WithChanLength(chanLength int) ClientOption {
	return func(config *ClientConfig) {
		config.chanLength = chanLength
	}
}

func WithReadDeadline(timeout time.Duration) ClientOption {
	return func(config *ClientConfig) {
		config.readDeadline = timeout
	}
}

func WithWriteDeadline(timeout time.Duration) ClientOption {
	return func(config *ClientConfig) {
		config.writeDeadline = timeout
	}
}

func WithReadLimit(readLimit int64) ClientOption {
	return func(config *ClientConfig) {
		config.readLimit = readLimit
	}
}

func WithSupportPing(period time.Duration) ClientOption {
	return func(config *ClientConfig) {
		if period > 0 {
			config.supportPing = true
			// 允许用户自定义周期，如果用户不传（即 period=0），则使用默认值
			config.pingPeriod = period
		} else {
			config.supportPing = false
		}
	}
}

func NewClient(conn *websocket.Conn, opts ...ClientOption) (*Client, error) {
	var defaultConnectConfig = &ClientConfig{
		chanLength:    1024,
		readDeadline:  10 * time.Second,
		writeDeadline: 10 * time.Second,
		readLimit:     0,
		supportPing:   false, // 默认关闭 Ping
		// 默认 Ping 周期应小于 ReadDeadline，通常为其 80%-90%
		pingPeriod: (10 * time.Second * 8) / 10,
	}

	for _, opt := range opts {
		opt(defaultConnectConfig)
	}

	client := &Client{
		conn:         conn,
		sendBuffer:   make(chan *Message, defaultConnectConfig.chanLength),
		ctx:          context.Background(),
		ClientConfig: defaultConnectConfig,
	}

	client.conn.SetPongHandler(func(appData string) error {
		return client.conn.SetReadDeadline(time.Now().Add(client.readDeadline))
	})

	if client.readLimit > 0 {
		conn.SetReadLimit(client.readLimit)
	}
	err := conn.SetReadDeadline(time.Now().Add(client.readDeadline))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Start() {
	go c.read()
	go c.write()
}

func (c *Client) read() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("read error:", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Close()
				return
			}

			if c.OnError != nil {
				c.OnError(c, err)
			}
			c.Close()

			return
		}

		err = c.conn.SetReadDeadline(time.Now().Add(c.readDeadline))
		if err != nil {
			if c.OnError != nil {
				c.OnError(c, err)
			}
			c.Close()
			return
		}

		if c.OnMessage != nil {
			c.OnMessage(c, message)
		}
	}
}

func (c *Client) writeMessage(messageType int, data []byte) bool {
	err := c.conn.SetWriteDeadline(time.Now().Add(c.writeDeadline))
	if err != nil {
		if c.OnError != nil {
			c.OnError(c, err)
		}
		c.Close()
		return false // 表示失败
	}

	err = c.conn.WriteMessage(messageType, data)
	if err != nil {
		if c.OnError != nil {
			c.OnError(c, err)
		}
		c.Close()
		return false // 表示失败
	}
	return true // 表示成功
}

func (c *Client) write() {
	if !c.supportPing {
		for msg := range c.sendBuffer {
			if !c.writeMessage(msg.Type, msg.Data) {
				// 发送失败，writeMessage 内部已经处理了关闭逻辑
				return
			}
		}
		return
	}

	// 支持 ping 的新逻辑
	ticker := time.NewTicker(c.pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.sendBuffer:
			if !ok {
				c.writeMessage(websocket.CloseMessage, []byte{})
				return
			}
			if !c.writeMessage(msg.Type, msg.Data) {
				return
			}

		case <-ticker.C:
			if !c.writeMessage(websocket.PingMessage, nil) {
				return
			}
		}
	}
}

func (c *Client) SendText(msg []byte) error {
	select {
	case c.sendBuffer <- &Message{Type: websocket.TextMessage, Data: msg}:
		return nil
	default:
		return fmt.Errorf("send buffer is full")
	}
}

func (c *Client) SendBinary(msg []byte) error {
	select {
	case c.sendBuffer <- &Message{Type: websocket.BinaryMessage, Data: msg}:
		return nil
	default:
		return fmt.Errorf("send buffer is full")
	}
}

// GetBaseClient 获取基础客户端实例
func (c *Client) GetBaseClient() *Client {
	return c
}

// SetContextValue 设置上下文值（非泛型版本）
func (c *Client) SetContextValue(key string, value interface{}) {
	c.ctxMu.Lock()
	defer c.ctxMu.Unlock()
	c.ctx = context.WithValue(c.ctx, key, value)
}

// GetContextValue 获取上下文值
func (c *Client) GetContextValue(key string) interface{} {
	c.ctxMu.Lock()
	defer c.ctxMu.Unlock()
	return c.ctx.Value(key)
}

// GetContextString 获取字符串类型的上下文值
func (c *Client) GetContextString(key string) string {
	if value := c.GetContextValue(key); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetContextInt 获取整数类型的上下文值
func (c *Client) GetContextInt(key string) int {
	if value := c.GetContextValue(key); value != nil {
		if i, ok := value.(int); ok {
			return i
		}
	}
	return 0
}

// GetContextBool 获取布尔类型的上下文值
func (c *Client) GetContextBool(key string) bool {
	if value := c.GetContextValue(key); value != nil {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return false
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		close(c.sendBuffer)

		err := c.conn.Close()
		if err != nil {
			if c.OnError != nil {
				c.OnError(c, err)
			}
			return
		}
		if c.OnClose != nil {
			c.OnClose(c)
		}
	})
}
