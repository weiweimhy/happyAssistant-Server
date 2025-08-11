package service

import (
	"errors"
	"fmt"
	"happyAssistant/internal/model"
	"happyAssistant/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
)

// UserService 用户服务
// 处理用户相关的业务逻辑，如登录、注册、用户信息管理等
type UserService struct {
	userRepo *repository.UserRepository
	labRepo  *repository.LabRepository
	roleRepo *repository.RoleRepository
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	repoManager := repository.GetRepositoryManager()
	return &UserService{
		userRepo: repoManager.GetUserRepository(),
		labRepo:  repoManager.GetLabRepository(),
		roleRepo: repoManager.GetRoleRepository(),
	}
}

// Login 用户登录
// 处理微信小程序登录，验证js_code并返回用户信息和实验室信息
func (us *UserService) Login(jsCode string) (*model.LoginResponse, error) {
	log.Infof("Processing login request with js_code: %s", jsCode)

	// 1. 验证js_code（这里需要调用微信API）
	openID, _, err := us.validateWechatCode(jsCode)
	if err != nil {
		return nil, fmt.Errorf("failed to validate wechat code: %w", err)
	}

	// 2. 查找或创建用户
	user, err := us.FindOrCreateUser(openID)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// 3. 获取用户默认实验室信息
	labInfo, err := us.getUserLabInfo(user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user lab info: %w", err)
	}

	// 4. 构建登录响应
	loginResp := &model.LoginResponse{
		User:    user,
		LabInfo: labInfo,
	}

	log.Infof("Login successful for user: %s", user.Id)
	return loginResp, nil
}

// validateWechatCode 验证微信小程序登录凭证
// 调用微信API验证js_code并获取openid和session_key
func (us *UserService) validateWechatCode(jsCode string) (string, string, error) {
	// TODO: 实现微信API调用
	// 这里需要调用微信小程序的登录API
	// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html

	// 临时返回模拟数据，实际项目中需要调用微信API
	if jsCode == "" {
		return "", "", errors.New("invalid js_code")
	}

	// 模拟微信API返回
	openID := "mock_openid_" + jsCode
	sessionKey := "mock_session_key_" + jsCode

	log.Infof("Validated wechat code, openID: %s", openID)
	return openID, sessionKey, nil
}

// findOrCreateUser 查找或创建用户
// 根据openID查找用户，如果不存在则创建新用户
func (us *UserService) FindOrCreateUser(openID string) (*model.User, error) {
	// 先尝试查找现有用户
	user, err := us.userRepo.FindByOpenID(openID)
	if err == nil && user != nil {
		log.Infof("Found existing user: %s", user.Id)
		return user, nil
	}

	// 如果用户不存在，创建新用户
	newUser := &model.User{
		Id:        generateUserID(),
		Name:      "新用户",
		Avatar:    "",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		// 注意：这里需要添加open_id字段到User model，或者使用其他方式存储
		// 暂时使用name字段存储openID作为临时方案
	}

	err = us.userRepo.Create(newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Infof("Created new user: %s", newUser.Id)
	return newUser, nil
}

// getUserLabInfo 获取用户实验室信息
// 获取用户当前选中的实验室信息和角色信息
func (us *UserService) getUserLabInfo(userID string) (*model.LoginLabInfo, error) {
	// 获取用户默认实验室（这里简化处理，实际可能需要从用户配置中获取）
	lab, err := us.labRepo.GetDefaultLab()
	if err != nil {
		return nil, fmt.Errorf("failed to get default lab: %w", err)
	}

	// 获取实验室的所有角色
	roles, err := us.roleRepo.GetRolesByLabID(lab.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get lab roles: %w", err)
	}

	// 获取用户在实验室中的角色（这里简化处理，实际可能需要查询用户-角色关联表）
	userRole := us.getDefaultUserRole(roles)
	if userRole == nil {
		return nil, errors.New("failed to get user role")
	}

	labInfo := &model.LoginLabInfo{
		Lab:        lab,
		Roles:      roles,
		UserRoleId: userRole.Id,
		UserRole:   userRole,
	}

	return labInfo, nil
}

// getDefaultUserRole 获取用户默认角色
// 这里简化处理，实际项目中需要根据用户权限和实验室配置来确定
func (us *UserService) getDefaultUserRole(roles []*model.Role) *model.Role {
	if len(roles) == 0 {
		return nil
	}

	// 优先返回学生角色，如果没有则返回第一个角色
	for _, role := range roles {
		if role.Name == "学生" {
			return role
		}
	}

	return roles[0]
}

// generateUserID 生成用户ID
// 生成唯一的用户标识符
func generateUserID() string {
	// TODO: 实现UUID生成或使用其他唯一ID生成策略
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}
