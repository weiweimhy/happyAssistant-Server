# 微信小程序后端框架

## 项目特性

- **WebSocket Hub 模式**: 基于 `gorilla/websocket` 构建的健壮 WebSocket 框架
- **Protocol Buffers**: 使用 protobuf 定义通信协议，确保数据格式统一高效
- **Context 管理**: 支持客户端上下文信息存储和访问
- **泛型数据库操作**: 使用 Go 1.23+ 泛型实现的通用 CRUD 操作
- **配置管理**: 支持 YAML 配置文件，支持 debug/release 环境切换
- **结构化日志**: 集成 logrus 日志库，支持日志级别和文件输出
- **MongoDB 支持**: 完整的 MongoDB 连接池和操作封装
- **接口驱动设计**: 通过接口实现业务逻辑与底层实现的解耦
- **协议控制器**: 统一的协议路由和处理机制
- **自动构建脚本**: 提供 protobuf 编译和工具安装脚本

## 项目结构

```
happyAssistant/
├── cmd/                    # 应用程序入口
│   └── server/
│       └── main.go        # 主程序入口
├── configs/               # 配置文件
│   ├── config_debug.yaml  # 调试环境配置
│   └── config_release.yaml # 生产环境配置
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   ├── controller/       # 控制器层
│   │   └── protocol_controller.go # 协议控制器
│   ├── customUtils/      # 自定义工具
│   ├── initialize/       # 初始化模块
│   ├── logger/           # 日志管理
│   ├── model/            # 数据模型 (protobuf生成)
│   │   ├── lab.pb.go
│   │   ├── permission.pb.go
│   │   ├── protocol.pb.go
│   │   ├── role.pb.go
│   │   └── user.pb.go
│   ├── repository/       # 数据访问层
│   │   ├── repository.go      # 泛型CRUD操作
│   │   ├── repository_manager.go
│   │   ├── user_repository.go
│   │   ├── lab_repository.go
│   │   └── role_repository.go
│   └── service/          # 业务逻辑层
│       ├── user_service.go
│       └── lab_service.go
├── pkg/                  # 公共包
│   └── wshub/           # WebSocket Hub
│       ├── websocket_hub.go
│       └── websocket_client.go
├── proto/               # Protocol Buffers 定义
│   ├── protocol.proto   # 通信协议
│   ├── user.proto       # 用户信息
│   ├── lab.proto        # 实验室信息
│   ├── role.proto       # 角色信息
│   └── permission.proto # 权限定义
├── go.mod               # Go 模块文件
├── go.sum               # 依赖校验和
├── install_protoc_tools.sh    # protoc工具安装脚本
├── install_protoc_tools.ps1   # Windows版本安装脚本
├── protocol_build.sh          # protobuf编译脚本
├── protocol_build.ps1         # Windows版本编译脚本
└── readme.md            # 项目文档
```

## 系统架构

### 分层架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                    WebSocket Hub Layer                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   Connection    │  │   Context       │  │   Event      │ │
│  │   Management    │  │   Management    │  │   Handling   │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                   Protocol Controller Layer                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   Protocol      │  │   Message       │  │   Response   │ │
│  │   Routing       │  │   Parsing       │  │   Building   │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                     Service Layer                           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   User          │  │   Lab           │  │   Role       │ │
│  │   Service       │  │   Service       │  │   Service    │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                   Repository Layer                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   Generic       │  │   User          │  │   Lab        │ │
│  │   CRUD          │  │   Repository    │  │   Repository │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                    MongoDB Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   Connection    │  │   Collection    │  │   Index      │ │
│  │   Pool          │  │   Management    │  │   Management │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 核心组件说明

#### 1. WebSocket Hub
- **单例模式**: 全局唯一的 WebSocket 服务实例
- **事件驱动**: 通过回调函数处理连接、消息、错误事件
- **Context 管理**: 支持客户端会话信息存储
- **工厂模式**: 支持自定义客户端创建

#### 2. Protocol Controller
- **协议路由**: 根据协议类型分发到相应处理器
- **消息解析**: 统一的 protobuf 消息解析机制
- **响应构建**: 标准化的响应格式和错误处理

#### 3. Service Layer
- **业务逻辑**: 封装核心业务处理逻辑
- **事务管理**: 确保数据一致性
- **权限控制**: 基于角色的访问控制

#### 4. Repository Layer
- **泛型CRUD**: 通用的数据库操作接口
- **数据访问**: 封装 MongoDB 操作细节
- **查询优化**: 支持复杂查询和索引优化

## Protocol Buffers 协议定义

### 核心协议文件

#### protocol.proto - 通信协议
```protobuf
// 协议类型枚举
enum ProtocolType {
  UNKNOWN = 0;        // 未知协议类型
  LOGIN_REQ = 1;      // 登录请求协议
  LOGIN_RESP = 2;     // 登录响应协议
}

// 响应状态码枚举
enum RESP_CODE {
  ERROR = 0;    // 操作失败
  SUCCESS = 1;  // 操作成功
}

// 基础请求协议
message BaseRequest {
  ProtocolType type = 1;  // 协议类型
  bytes data = 2;         // 具体请求数据
}

// 基础响应协议
message BaseResponse {
  ProtocolType type = 1;  // 协议类型
  RESP_CODE result = 2;   // 操作结果
  string msg = 3;         // 结果消息
  bytes data = 4;         // 具体响应数据
  int64 timestamp = 5;    // 时间戳
}

// 登录请求协议
message LoginRequest {
  string js_code = 1;     // 微信小程序登录凭证
}

// 登录响应协议
message LoginResponse {
  user.User user = 1;           // 用户信息
  LoginLabInfo labInfo = 2;     // 实验室信息
}
```

#### user.proto - 用户信息
```protobuf
message User {
  string id = 1;              // 用户唯一标识符
  string name = 2;            // 用户姓名
  string avatar = 3;          // 用户头像URL
  string phone_number = 4;    // 用户手机号码
  string email = 5;           // 用户邮箱地址
  int64 created_at = 6;       // 用户创建时间戳
  int64 updated_at = 7;       // 用户信息更新时间戳
  repeated string lib_ids = 8; // 用户所属的实验室ID列表
}
```

#### lab.proto - 实验室信息
```protobuf
message Lab {
  string id = 1;                           // 实验室唯一标识符
  string name = 2;                         // 实验室名称
  string desc = 3;                         // 实验室描述信息
  string create_id = 4;                    // 实验室创建者ID
  string owner_id = 5;                     // 实验室所有者ID
  int64 create_at = 6;                     // 实验室创建时间戳
  int64 update_at = 7;                     // 实验室信息更新时间戳
  repeated string role_ids = 8;            // 实验室中定义的角色ID列表
  map<string, string> user_role_map = 9;   // 用户角色映射表
  repeated role.Role roles = 10;           // 实验室中定义的完整角色信息列表
}
```

#### role.proto - 角色信息
```protobuf
message Role {
  string id = 1;                // 角色唯一标识符
  string name = 2;              // 角色名称
  uint64 permission_flags = 3;  // 权限标志位
}
```

#### permission.proto - 权限枚举
```protobuf
enum Permission {
  UNKNOWN = 0;        // 未知权限
  ORDER_CREATE = 1;   // 创建订单权限 (1 << 0)
  ORDER_UPDATE = 2;   // 更新订单权限 (1 << 1)
  ORDER_DELETE = 4;   // 删除订单权限 (1 << 2)
}
```

### 协议编译

#### 自动安装工具
```bash
# Linux/macOS
chmod +x install_protoc_tools.sh
./install_protoc_tools.sh

# Windows PowerShell
.\install_protoc_tools.ps1
```

#### 编译协议文件
```bash
# Linux/macOS
chmod +x protocol_build.sh
./protocol_build.sh

# Windows PowerShell
.\protocol_build.ps1
```

#### 手动编译命令
```bash
# 编译所有proto文件
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/*.proto

# 注入BSON标签
protoc-go-inject-tag -input=internal/model/*.pb.go
```

#### 生成的文件
编译后会在`internal/model/`目录下生成以下Go文件：
- `protocol.pb.go` - 协议相关结构体
- `user.pb.go` - 用户相关结构体
- `lab.pb.go` - 实验室相关结构体
- `role.pb.go` - 角色相关结构体
- `permission.pb.go` - 权限相关结构体

## WebSocket Hub 框架

### 核心特性

- **接口驱动**: 通过 `IClient` 接口实现业务客户端与底层网络客户端的解耦
- **Context 管理**: 支持客户端上下文信息存储和访问
- **客户端工厂**: 支持通过工厂模式创建自定义的业务客户端（可选）
- **双向心跳**: 内置 Ping/Pong 处理逻辑，确保连接稳定
- **优雅的回调机制**: 通过事件回调处理所有网络事件
- **配置灵活**: 支持函数式选项模式自定义行为
- **类型安全**: 支持文本和二进制消息发送

### IClient 接口

```go
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
```

### Context 功能

#### 设置和获取上下文信息
```go
// 设置上下文信息
client.SetContextValue("user_id", "user_123")
client.SetContextValue("username", "张三")
client.SetContextValue("lab_id", "lab_456")
client.SetContextValue("is_online", true)

// 获取上下文信息
userID := client.GetContextString("user_id")
username := client.GetContextString("username")
labID := client.GetContextString("lab_id")
isOnline := client.GetContextBool("is_online")
```

### 快速上手

#### 1. 使用默认客户端

项目提供了默认的 `Client` 实现，可以直接使用，无需创建自定义客户端：

```go
// 直接使用默认的 wshub.Client
// 它已经实现了 IClient 接口，包含所有必要的方法
```

#### 2. 启动 WebSocket 服务

```go
// file: cmd/server/main.go
package main

import (
    "log"
    "time"
    "happyAssistant/pkg/wshub"
)

func main() {
    // 获取Hub单例
    hub := wshub.GetInstance()
    
    // 配置客户端选项
    hub.SetClientOptions(
        wshub.WithSupportPing(45*time.Second),
        wshub.WithChanLength(2048),
    )
    
    // 设置事件回调
    hub.OnOpen = func(client wshub.IClient) {
        log.Printf("客户端连接: %p", client.GetBaseClient())
        
        // 设置上下文信息
        client.SetContextValue("connection_time", time.Now().Unix())
        client.SetContextValue("user_id", "user_123")
        client.SetContextValue("username", "张三")
    }
    
    hub.OnMessage = func(client wshub.IClient, msg []byte) {
        // 从上下文获取信息
        connectionTime := client.GetContextInt("connection_time")
        userID := client.GetContextString("user_id")
        username := client.GetContextString("username")
        
        log.Printf("收到来自用户 %s(%s) 的消息，连接时间: %d", username, userID, connectionTime)
        
        // 处理消息
        handleMessage(client, msg)
    }
    
    hub.OnClose = func(client wshub.IClient) {
        userID := client.GetContextString("user_id")
        log.Printf("客户端断开: UserID=%s", userID)
    }
    
    hub.OnError = func(client wshub.IClient, err error) {
        log.Printf("发生错误: %v", err)
    }
    
    // 启动服务
    log.Println("WebSocket Hub 启动，端口 8080...")
    if err := hub.Start("/ws", 8080); err != nil {
        log.Fatalf("启动失败: %v", err)
    }
}

func handleMessage(client wshub.IClient, msg []byte) {
    // 处理消息逻辑
    client.SendText([]byte("收到消息: " + string(msg)))
}
```

### 协议处理示例

#### 客户端发送登录请求
```go
import (
    "happyAssistant/internal/model"
    "google.golang.org/protobuf/proto"
)

// 创建登录请求
loginReq := &model.LoginRequest{
    JsCode: "wx_login_code_123",
}

// 序列化登录请求
loginData, err := proto.Marshal(loginReq)
if err != nil {
    log.Fatal("序列化登录请求失败:", err)
}

// 创建基础请求
baseReq := &model.BaseRequest{
    Type: model.ProtocolType_LOGIN_REQ,
    Data: loginData,
}

// 序列化基础请求
requestData, err := proto.Marshal(baseReq)
if err != nil {
    log.Fatal("序列化基础请求失败:", err)
}

// 通过WebSocket发送
client.SendBinary(requestData)
```

#### 服务器处理登录请求
```go
hub.OnMessage = func(client IClient, msg []byte) {
    // 解析基础请求
    var baseReq model.BaseRequest
    err := proto.Unmarshal(msg, &baseReq)
    if err != nil {
        log.Printf("解析基础请求失败: %v", err)
        return
    }
    
    switch baseReq.Type {
    case model.ProtocolType_LOGIN_REQ:
        // 解析登录请求
        var loginReq model.LoginRequest
        err := proto.Unmarshal(baseReq.Data, &loginReq)
        if err != nil {
            log.Printf("解析登录请求失败: %v", err)
            return
        }
        
        // 处理登录逻辑
        user, labInfo, err := handleLogin(loginReq.JsCode)
        if err != nil {
            sendErrorResponse(client, err)
            return
        }
        
        // 创建登录响应
        loginResp := &model.LoginResponse{
            User:    user,
            LabInfo: labInfo,
        }
        
        // 发送登录响应
        sendLoginResponse(client, loginResp)
        
    default:
        log.Printf("未知协议类型: %v", baseReq.Type)
    }
}
```

## 配置管理

### 配置文件结构

项目支持 YAML 格式的配置文件，支持 debug 和 release 两种模式：

```yaml
# 服务端配置
server:
  port: 8080
  route: "/ws"

# MongoDB配置
mongodb:
  uri: "mongodb://localhost:27017"
  database: "lzdb_debug"  # debug模式使用测试数据库
  username: ""
  password: ""
  timeout: 10s
  opTimeout: 5s

# 日志配置
log:
  level: info  # debug, info, warn, error
  file: ""     # 留空输出到控制台，指定文件路径输出到文件
```

### 环境配置对比

| 配置项 | Debug 模式 | Release 模式 |
|--------|------------|--------------|
| 端口 | 8080 | 9300 |
| 路由 | /ws | /wss |
| 数据库 | lzdb_debug | lzdb_release |
| 日志级别 | info | error |
| 日志输出 | 控制台 | 文件 |

### 配置加载

```go
// 在 main.go 中加载配置
config.LoadConfig("configs/config_debug.yaml")

// 使用配置
port := config.Cfg.Server.Port
dbName := config.Cfg.MongoDB.Database
logLevel := config.Cfg.Log.Level
```

## 泛型数据库操作

### 接口定义

所有数据模型都需要实现 `proto.Message` 接口（protobuf 自动生成）：

```go
// protobuf 自动生成的接口
type Message interface {
    ProtoReflect() protoreflect.Message
    Reset()
    String() string
    ProtoMessage()
}
```

### 数据模型示例

```go
// protobuf 生成的 User 结构体
type User struct {
    Id          string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
    Name        string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
    Avatar      string   `protobuf:"bytes,3,opt,name=avatar,proto3" json:"avatar,omitempty"`
    PhoneNumber string   `protobuf:"bytes,4,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty"`
    Email       string   `protobuf:"bytes,5,opt,name=email,proto3" json:"email,omitempty"`
    CreatedAt   int64    `protobuf:"varint,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
    UpdatedAt   int64    `protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
    LibIds      []string `protobuf:"bytes,8,rep,name=lib_ids,json=libIds,proto3" json:"lib_ids,omitempty"`
}
```

### 泛型 CRUD 操作

```go
// 插入操作
user := &model.User{
    Id: "123", 
    Name: "张三",
    CreatedAt: time.Now().Unix(),
}
err := repository.InsertOne(userCollection, user)

// 查询操作
filter := bson.M{"name": "张三"}
user, err := repository.FindOne[model.User](userCollection, filter)
users, err := repository.FindMany[model.User](userCollection, filter)

// 更新操作
update := bson.M{"$set": bson.M{
    "email": "new@example.com",
    "updated_at": time.Now().Unix(),
}}
err := repository.UpdateOne[model.User](userCollection, filter, update)

// 删除操作
err := repository.DeleteOne(userCollection, filter)

// 统计操作
count, err := repository.Count[model.User](userCollection, filter)
```

### 支持的数据库操作

- `InsertOne`: 插入单个文档
- `InsertMany`: 批量插入文档
- `FindOne`: 查询单个文档
- `FindMany`: 查询多个文档
- `UpdateOne`: 更新单个文档
- `UpdateMany`: 批量更新文档
- `ReplaceOne`: 替换单个文档
- `DeleteOne`: 删除单个文档
- `DeleteMany`: 批量删除文档
- `Count`: 统计文档数量

## 日志系统

### 集成 logrus

项目使用 logrus 作为日志库，支持结构化日志和多种输出格式：

```go
import "github.com/sirupsen/logrus"

// 设置日志级别
logrus.SetLevel(logrus.InfoLevel)

// 结构化日志
logrus.WithFields(logrus.Fields{
    "user_id": "123",
    "action": "login",
    "lab_id": "lab_456",
}).Info("用户登录成功")
```

### 日志配置

```go
func InitLogger() {
    // 设置日志级别
    level, err := logrus.ParseLevel(config.Cfg.Log.Level)
    if err != nil {
        level = logrus.InfoLevel
    }
    logrus.SetLevel(level)

    // 设置日志格式
    logrus.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
    })

    // 设置日志输出文件
    if config.Cfg.Log.File != "" {
        f, err := os.OpenFile(config.Cfg.Log.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
        if err != nil {
            logrus.Fatalf("open log file error: %v", err)
        }
        logrus.SetOutput(f)
    }
}
```

## MongoDB 连接管理

### 连接初始化

```go
// 初始化 MongoDB 客户端
mongoClient := initialize.InitMongoDBClient(&config.Cfg.MongoDB)

// 获取集合
userCollection := mongoClient.Collection(config.Cfg.MongoDB.Database, "users")
labCollection := mongoClient.Collection(config.Cfg.MongoDB.Database, "labs")
roleCollection := mongoClient.Collection(config.Cfg.MongoDB.Database, "roles")
```

### 连接池配置

- 支持用户名密码认证
- 可配置最大连接池大小
- 支持连接超时和操作超时设置
- 全局单例模式，线程安全

## API 文档

### WebSocket 连接

#### 连接地址
- **Debug 模式**: `ws://localhost:8080/ws`
- **Release 模式**: `wss://your-domain.com/wss`

#### 消息格式
所有消息都使用 Protocol Buffers 二进制格式，包含在 `BaseRequest` 和 `BaseResponse` 中。

### 协议列表

#### 1. 登录协议

**请求**: `LOGIN_REQ`
```protobuf
message LoginRequest {
    string js_code = 1;  // 微信小程序登录凭证
}
```

**响应**: `LOGIN_RESP`
```protobuf
message LoginResponse {
    user.User user = 1;           // 用户信息
    LoginLabInfo labInfo = 2;     // 实验室信息
}
```

**使用示例**:
```javascript
// 客户端 JavaScript 示例
const ws = new WebSocket('ws://localhost:8080/ws');

// 发送登录请求
const loginRequest = {
    type: 1,  // LOGIN_REQ
    data: new Uint8Array([...])  // 序列化的 LoginRequest
};

ws.send(new Uint8Array([...]));  // 序列化的 BaseRequest
```

### 错误处理

所有错误响应都遵循统一的格式：

```protobuf
message BaseResponse {
    ProtocolType type = 1;     // 对应的请求协议类型
    RESP_CODE result = 2;      // ERROR = 0
    string msg = 3;            // 错误描述
    bytes data = 4;            // 空或错误详情
    int64 timestamp = 5;       // 时间戳
}
```

## 部署指南

### 开发环境部署

#### 1. 环境要求

- **Go**: 1.23+ (支持泛型)
- **MongoDB**: 4.0+
- **Protocol Buffers**: 3.0+

#### 2. 安装依赖

```bash
# 安装Go依赖
go mod tidy

# 安装Protocol Buffers工具
chmod +x install_protoc_tools.sh
./install_protoc_tools.sh
```

#### 3. 编译协议文件

```bash
chmod +x protocol_build.sh
./protocol_build.sh
```

#### 4. 配置数据库

确保 MongoDB 服务正在运行，并修改 `configs/config_debug.yaml` 中的数据库连接信息。

#### 5. 运行服务

```bash
go run cmd/server/main.go
```

### 生产环境部署

#### 1. 构建二进制文件

```bash
# 构建可执行文件
go build -o server cmd/server/main.go

# 或使用交叉编译
GOOS=linux GOARCH=amd64 go build -o server cmd/server/main.go
```

#### 2. 配置文件

创建生产环境配置文件 `configs/config_release.yaml`：

```yaml
server:
  port: 9300
  route: "/wss"

mongodb:
  uri: "mongodb://your-mongodb-server:27017"
  database: "lzdb_release"
  username: "your-username"
  password: "your-password"
  timeout: 10s
  opTimeout: 5s

log:
  level: error
  file: "/var/log/server.log"
```

#### 3. 使用 systemd 管理服务

创建服务文件 `/etc/systemd/system/websocket-server.service`：

```ini
[Unit]
Description=WebSocket Server
After=network.target

[Service]
Type=simple
User=websocket
WorkingDirectory=/opt/websocket-server
ExecStart=/opt/websocket-server/server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable websocket-server
sudo systemctl start websocket-server
```

#### 4. 使用 Docker 部署

创建 `Dockerfile`：

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs

EXPOSE 9300
CMD ["./server"]
```

构建和运行：
```bash
docker build -t websocket-server .
docker run -d -p 9300:9300 --name websocket-server websocket-server
```

### 反向代理配置

#### Nginx 配置示例

```nginx
upstream websocket_backend {
    server 127.0.0.1:9300;
}

server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location /wss {
        proxy_pass http://websocket_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 86400;
    }
}
```

## 监控和运维

### 健康检查

添加健康检查端点：

```go
func healthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now().Unix(),
        "version": "1.0.0",
    })
}

// 在 main.go 中添加
http.HandleFunc("/health", healthCheck)
```

### 性能监控

#### 连接数监控

```go
type Metrics struct {
    ActiveConnections int64
    TotalMessages     int64
    ErrorCount        int64
}

var metrics = &Metrics{}

// 在 WebSocket 事件中更新指标
hub.OnOpen = func(client wshub.IClient) {
    atomic.AddInt64(&metrics.ActiveConnections, 1)
}

hub.OnClose = func(client wshub.IClient) {
    atomic.AddInt64(&metrics.ActiveConnections, -1)
}
```

#### 内存使用监控

```go
func getMemoryStats() map[string]interface{} {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return map[string]interface{}{
        "alloc":      m.Alloc,
        "total_alloc": m.TotalAlloc,
        "sys":        m.Sys,
        "num_gc":     m.NumGC,
    }
}
```

### 日志轮转

使用 logrotate 配置日志轮转：

```bash
# /etc/logrotate.d/websocket-server
/var/log/server.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 644 websocket websocket
    postrotate
        systemctl reload websocket-server
    endscript
}
```

## 测试指南

### 单元测试

#### 测试 Repository 层

```go
// repository_test.go
func TestUserRepository_FindByID(t *testing.T) {
    // 设置测试数据库
    client := setupTestDB(t)
    collection := client.Collection("test_db", "users")
    
    // 创建测试数据
    user := &model.User{
        Id: "test_user_1",
        Name: "测试用户",
    }
    err := InsertOne(collection, user)
    require.NoError(t, err)
    
    // 测试查询
    foundUser, err := FindOne[model.User](collection, bson.M{"_id": "test_user_1"})
    require.NoError(t, err)
    require.Equal(t, "测试用户", foundUser.Name)
}
```

#### 测试 Service 层

```go
// service_test.go
func TestUserService_Login(t *testing.T) {
    // 创建 mock repository
    mockRepo := &MockUserRepository{}
    
    // 设置期望行为
    mockRepo.On("FindByID", "test_user").Return(&model.User{
        Id: "test_user",
        Name: "测试用户",
    }, nil)
    
    // 创建 service 实例
    service := &UserService{
        userRepo: mockRepo,
    }
    
    // 执行测试
    user, err := service.Login("test_js_code")
    require.NoError(t, err)
    require.Equal(t, "测试用户", user.Name)
    
    // 验证 mock 调用
    mockRepo.AssertExpectations(t)
}
```

### 集成测试

#### WebSocket 连接测试

```go
func TestWebSocketConnection(t *testing.T) {
    // 启动测试服务器
    go func() {
        hub := wshub.GetInstance()
        hub.Start("/ws", 8081)
    }()
    
    time.Sleep(100 * time.Millisecond)
    
    // 连接 WebSocket
    conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081/ws", nil)
    require.NoError(t, err)
    defer conn.Close()
    
    // 发送测试消息
    testMsg := []byte("test message")
    err = conn.WriteMessage(websocket.TextMessage, testMsg)
    require.NoError(t, err)
    
    // 接收响应
    _, msg, err := conn.ReadMessage()
    require.NoError(t, err)
    require.Contains(t, string(msg), "test message")
}
```

### 性能测试

#### 并发连接测试

```go
func BenchmarkWebSocketConnections(b *testing.B) {
    // 启动服务器
    go func() {
        hub := wshub.GetInstance()
        hub.Start("/ws", 8082)
    }()
    
    time.Sleep(100 * time.Millisecond)
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8082/ws", nil)
            if err != nil {
                b.Fatal(err)
            }
            conn.Close()
        }
    })
}
```

## 故障排除

### 常见问题

#### 1. Protocol Buffers 编译失败

**问题**: `protoc: command not found`
**解决方案**:
```bash
# 安装 protoc
# Ubuntu/Debian
sudo apt install protobuf-compiler

# macOS
brew install protobuf

# Windows
# 下载 protoc 二进制文件并添加到 PATH
```

#### 2. MongoDB 连接失败

**问题**: `connection refused`
**解决方案**:
```bash
# 检查 MongoDB 服务状态
sudo systemctl status mongod

# 启动 MongoDB 服务
sudo systemctl start mongod

# 检查连接配置
mongo --host localhost --port 27017
```

#### 3. WebSocket 连接被拒绝

**问题**: `websocket: bad handshake`
**解决方案**:
- 检查服务器端口是否正确
- 确认防火墙设置
- 验证 WebSocket 路由配置

#### 4. 内存泄漏

**问题**: 内存使用持续增长
**解决方案**:
```go
// 定期清理断开的连接
func cleanupDisconnectedClients() {
    ticker := time.NewTicker(5 * time.Minute)
    for range ticker.C {
        // 清理逻辑
    }
}
```

### 调试技巧

#### 1. 启用详细日志

```go
// 设置日志级别为 debug
logrus.SetLevel(logrus.DebugLevel)
```

#### 2. 使用 pprof 进行性能分析

```go
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    // ... 其他代码
}
```

#### 3. 监控 Goroutine 数量

```go
func monitorGoroutines() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        logrus.Infof("Goroutines: %d", runtime.NumGoroutine())
    }
}
```

## 项目最佳实践

### 1. 代码组织

- **分层架构**: 遵循 MVC 模式，分离业务逻辑和数据访问
- **接口驱动**: 通过接口实现依赖倒置，提高可测试性
- **错误处理**: 统一的错误处理机制，避免错误传播
- **日志记录**: 结构化日志，便于问题排查

### 2. 数据库设计

- **文档设计**: 合理设计 MongoDB 文档结构，避免过度嵌套
- **索引优化**: 为常用查询字段创建索引
- **连接池**: 合理配置连接池大小，避免连接泄漏

### 3. WebSocket 使用

- **心跳机制**: 使用内置的 Ping/Pong 机制保持连接
- **Context 管理**: 合理使用 Context 存储会话信息
- **错误处理**: 实现完整的错误处理和重连机制
- **消息格式**: 使用 Protocol Buffers 确保消息格式统一

### 4. 配置管理

- **环境分离**: 严格区分开发、测试、生产环境配置
- **敏感信息**: 避免在代码中硬编码敏感信息
- **配置验证**: 启动时验证配置的完整性和正确性

### 5. 安全考虑

- **输入验证**: 验证所有客户端输入
- **权限控制**: 实现基于角色的访问控制
- **数据加密**: 敏感数据传输使用加密
- **速率限制**: 防止恶意请求和 DoS 攻击

## 依赖说明

### 核心依赖

- **Go**: 1.23+ (支持泛型)
- **MongoDB**: 4.0+
- **gorilla/websocket**: WebSocket 实现
- **logrus**: 结构化日志
- **gopkg.in/yaml.v3**: YAML 配置解析
- **go.mongodb.org/mongo-driver/v2**: MongoDB 驱动
- **google.golang.org/protobuf**: Protocol Buffers

### 开发依赖

- **protoc**: Protocol Buffers 编译器
- **protoc-gen-go**: Go 语言 protobuf 插件
- **protoc-go-inject-tag**: BSON 标签注入工具

### 版本兼容性

| 组件 | 最低版本 | 推荐版本 |
|------|----------|----------|
| Go | 1.23 | 1.24+ |
| MongoDB | 4.0 | 6.0+ |
| Protocol Buffers | 3.0 | 3.21+ |

## 贡献指南

### 开发流程

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释和文档
- 编写单元测试

### 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
feat: 添加新功能
fix: 修复bug
docs: 更新文档
style: 代码格式调整
refactor: 代码重构
test: 添加测试
chore: 构建过程或辅助工具的变动
```

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 实现基础 WebSocket Hub 框架
- 支持 Protocol Buffers 协议
- 集成 MongoDB 数据库操作
- 提供泛型 CRUD 操作接口

### 计划功能
- [ ] 支持 Redis 缓存
- [ ] 添加 JWT 认证
- [ ] 实现消息队列
- [ ] 支持集群部署
- [ ] 添加 Prometheus 监控

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue
- 发送邮件
- 参与讨论

---

**注意**: 这是一个持续开发的项目，文档会随着功能更新而更新。请定期查看最新版本。