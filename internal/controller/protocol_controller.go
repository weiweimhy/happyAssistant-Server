package controller

import (
	"happyAssistant/internal/model"
	"happyAssistant/internal/service"
	"happyAssistant/pkg/wshub"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

// ProtocolController 协议控制器
// 负责解析客户端请求协议，并根据协议类型路由到相应的业务处理器
type ProtocolController struct {
	userService *service.UserService
	// 可以添加其他服务
}

// NewProtocolController 创建协议控制器实例
func NewProtocolController() *ProtocolController {
	return &ProtocolController{
		userService: service.NewUserService(),
	}
}

// HandleMessage 处理客户端消息
// 解析基础请求协议，并根据协议类型分发到相应的处理器
func (pc *ProtocolController) HandleMessage(client wshub.IClient, msg []byte) {
	// 解析基础请求协议
	var baseReq model.BaseRequest
	if err := proto.Unmarshal(msg, &baseReq); err != nil {
		log.Errorf("Failed to unmarshal base request: %v", err)
		pc.sendErrorResponse(client, model.ProtocolType_UNKNOWN, "Invalid request format")
		return
	}

	// 根据协议类型路由到相应的处理器
	switch baseReq.Type {
	case model.ProtocolType_LOGIN_REQ:
		pc.handleLoginRequest(client, baseReq.Data)
	default:
		log.Warnf("Unknown protocol type: %v", baseReq.Type)
		pc.sendErrorResponse(client, baseReq.Type, "Unknown protocol type")
	}
}

// handleLoginRequest 处理登录请求
func (pc *ProtocolController) handleLoginRequest(client wshub.IClient, data []byte) {
	// 解析登录请求
	var loginReq model.LoginRequest
	if err := proto.Unmarshal(data, &loginReq); err != nil {
		log.Errorf("Failed to unmarshal login request: %v", err)
		pc.sendErrorResponse(client, model.ProtocolType_LOGIN_REQ, "Invalid login request format")
		return
	}

	// 调用业务服务处理登录
	loginResp, err := pc.userService.Login(loginReq.JsCode)
	if err != nil {
		log.Errorf("Login failed: %v", err)
		pc.sendErrorResponse(client, model.ProtocolType_LOGIN_REQ, err.Error())
		return
	}

	// 发送成功响应
	pc.sendSuccessResponse(client, model.ProtocolType_LOGIN_RESP, loginResp)
}

// sendSuccessResponse 发送成功响应
func (pc *ProtocolController) sendSuccessResponse(client wshub.IClient, protocolType model.ProtocolType, data proto.Message) {
	// 序列化响应数据
	dataBytes, err := proto.Marshal(data)
	if err != nil {
		log.Errorf("Failed to marshal response data: %v", err)
		return
	}

	// 构建基础响应
	baseResp := &model.BaseResponse{
		Type:      protocolType,
		Result:    model.RESP_CODE_SUCCESS,
		Msg:       "Success",
		Data:      dataBytes,
		Timestamp: getCurrentTimestamp(),
	}

	// 序列化基础响应
	respBytes, err := proto.Marshal(baseResp)
	if err != nil {
		log.Errorf("Failed to marshal base response: %v", err)
		return
	}

	// 发送响应
	if err := client.SendBinary(respBytes); err != nil {
		log.Errorf("Failed to send response: %v", err)
	}
}

// sendErrorResponse 发送错误响应
func (pc *ProtocolController) sendErrorResponse(client wshub.IClient, protocolType model.ProtocolType, errorMsg string) {
	baseResp := &model.BaseResponse{
		Type:      protocolType,
		Result:    model.RESP_CODE_ERROR,
		Msg:       errorMsg,
		Data:      nil,
		Timestamp: getCurrentTimestamp(),
	}

	respBytes, err := proto.Marshal(baseResp)
	if err != nil {
		log.Errorf("Failed to marshal error response: %v", err)
		return
	}

	if err := client.SendBinary(respBytes); err != nil {
		log.Errorf("Failed to send error response: %v", err)
	}
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}
