package wshub

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	instance *WebSocketHub
	once     sync.Once
)

type ServerConfig struct {
	UpGrader websocket.Upgrader
}

type ServerOption func(*ServerConfig)

type WebSocketHub struct {
	OnOpen        func(client IClient)
	OnClose       func(client IClient)
	OnMessage     func(client IClient, msg []byte)
	OnError       func(client IClient, err error)
	clientFactory func(baseClient *Client) IClient
	*ServerConfig
	clientOptions []ClientOption
}

func getDefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		UpGrader: websocket.Upgrader{
			ReadBufferSize:  1024, // 默认读缓冲区大小
			WriteBufferSize: 1024, // 默认写缓冲区大小
			CheckOrigin: func(r *http.Request) bool {
				return true // 默认允许所有跨域请求
			},
		},
	}
}

func GetInstance() *WebSocketHub {
	once.Do(func() {
		instance = &WebSocketHub{
			ServerConfig: getDefaultServerConfig(),
		}
	})
	return instance
}

// WithReadBufferSize 设置读缓冲区大小
func WithReadBufferSize(size int) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.UpGrader.ReadBufferSize = size
	}
}

// WithWriteBufferSize 设置写缓冲区大小
func WithWriteBufferSize(size int) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.UpGrader.WriteBufferSize = size
	}
}

// WithCheckOrigin 设置跨域检查函数
func WithCheckOrigin(checkOrigin func(*http.Request) bool) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.UpGrader.CheckOrigin = checkOrigin
	}
}

// SetClientFactory 设置客户端工厂函数
func (wsh *WebSocketHub) SetClientFactory(factory func(baseClient *Client) IClient) {
	wsh.clientFactory = factory
}

// SetClientOptions 设置客户端选项
func (wsh *WebSocketHub) SetClientOptions(opts ...ClientOption) {
	wsh.clientOptions = opts
}

func (wsh *WebSocketHub) processRequest(w http.ResponseWriter, r *http.Request) {
	conn, err := wsh.UpGrader.Upgrade(w, r, nil)
	if err != nil {
		if wsh.OnError != nil {
			wsh.OnError(nil, err)
		}
		return
	}

	baseClient, err := NewClient(conn, wsh.clientOptions...)
	if err != nil {
		if wsh.OnError != nil {
			wsh.OnError(nil, err)
		}
		return
	}

	var client IClient
	if wsh.clientFactory != nil {
		client = wsh.clientFactory(baseClient)
	} else {
		client = baseClient
	}

	baseClient.OnMessage = func(_ IClient, msg []byte) {
		if wsh.OnMessage != nil {
			wsh.OnMessage(client, msg)
		}
	}
	baseClient.OnError = func(_ IClient, err error) {
		if wsh.OnError != nil {
			wsh.OnError(client, err)
		}
	}
	baseClient.OnClose = func(_ IClient) {
		if wsh.OnClose != nil {
			wsh.OnClose(client)
		}
	}

	if wsh.OnOpen != nil {
		wsh.OnOpen(client)
	}

	baseClient.Start()
}

func (wsh *WebSocketHub) Start(route string, port int, opts ...ServerOption) error {
	for _, opt := range opts {
		opt(wsh.ServerConfig)
	}

	http.HandleFunc(route, wsh.processRequest)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return err
}
