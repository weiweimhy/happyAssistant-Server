package repository

import (
	"sync"
)

// RepositoryManager 仓库管理器，使用单例模式管理所有 Repository 实例
type RepositoryManager struct {
	userRepo *UserRepository
	labRepo  *LabRepository
	roleRepo *RoleRepository
}

var (
	repositoryManager *RepositoryManager
	once              sync.Once
)

// GetRepositoryManager 获取仓库管理器单例
func GetRepositoryManager() *RepositoryManager {
	once.Do(func() {
		repositoryManager = &RepositoryManager{
			userRepo: NewUserRepository(),
			labRepo:  NewLabRepository(),
			roleRepo: NewRoleRepository(),
		}
	})
	return repositoryManager
}

// GetUserRepository 获取用户仓库
func (rm *RepositoryManager) GetUserRepository() *UserRepository {
	return rm.userRepo
}

// GetLabRepository 获取实验室仓库
func (rm *RepositoryManager) GetLabRepository() *LabRepository {
	return rm.labRepo
}

// GetRoleRepository 获取角色仓库
func (rm *RepositoryManager) GetRoleRepository() *RoleRepository {
	return rm.roleRepo
}
