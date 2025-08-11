package service

import (
	"fmt"
	"happyAssistant/internal/model"
	"happyAssistant/internal/repository"

	log "github.com/sirupsen/logrus"
)

// LabService 实验室服务
// 处理实验室相关的业务逻辑，如实验室管理、用户权限等
type LabService struct {
	labRepo  *repository.LabRepository
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

// NewLabService 创建实验室服务实例
func NewLabService() *LabService {
	repoManager := repository.GetRepositoryManager()
	return &LabService{
		labRepo:  repoManager.GetLabRepository(),
		userRepo: repoManager.GetUserRepository(),
		roleRepo: repoManager.GetRoleRepository(),
	}
}

// GetLabWithUsers 获取实验室及其用户信息
func (ls *LabService) GetLabWithUsers(labID string) (*model.Lab, []*model.User, error) {
	log.Infof("Getting lab with users: %s", labID)

	// 获取实验室信息
	lab, err := ls.labRepo.FindByID(labID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get lab: %w", err)
	}

	// 获取实验室的用户列表（这里需要根据实际业务逻辑实现）
	// 假设有一个方法可以获取实验室的用户
	users, err := ls.getLabUsers(labID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get lab users: %w", err)
	}

	return lab, users, nil
}

// CreateLab 创建实验室
func (ls *LabService) CreateLab(lab *model.Lab) error {
	log.Infof("Creating lab: %s", lab.Id)
	return ls.labRepo.Create(lab)
}

// GetLabRoles 获取实验室的角色列表
func (ls *LabService) GetLabRoles(labID string) ([]*model.Role, error) {
	log.Infof("Getting lab roles: %s", labID)
	return ls.roleRepo.GetRolesByLabID(labID)
}

// getLabUsers 获取实验室的用户列表
// 这是一个示例方法，实际实现需要根据数据库设计来调整
func (ls *LabService) getLabUsers(labID string) ([]*model.User, error) {
	// TODO: 实现获取实验室用户的逻辑
	// 这里可能需要查询用户-实验室关联表
	log.Infof("Getting users for lab: %s", labID)
	return []*model.User{}, nil
}
