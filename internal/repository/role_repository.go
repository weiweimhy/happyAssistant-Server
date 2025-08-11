package repository

import (
	"happyAssistant/internal/initialize"
	"happyAssistant/internal/model"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// RoleRepository 角色数据访问层
type RoleRepository struct {
	collection *mongo.Collection
}

// NewRoleRepository 创建角色仓库实例
func NewRoleRepository() *RoleRepository {
	client := initialize.GetMongoClient()
	collection := client.Collection("roles")
	return &RoleRepository{
		collection: collection,
	}
}

// GetRolesByLabID 根据实验室ID获取角色列表
func (rr *RoleRepository) GetRolesByLabID(labID string) ([]*model.Role, error) {
	log.Infof("Getting roles by lab ID: %s", labID)
	filter := bson.M{"lab_id": labID}
	results, err := FindMany[*model.Role](rr.collection, filter)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// FindByID 根据ID查找角色
func (rr *RoleRepository) FindByID(roleID string) (*model.Role, error) {
	log.Infof("Finding role by ID: %s", roleID)
	filter := bson.M{"_id": roleID}
	result, err := FindOne[*model.Role](rr.collection, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Create 创建角色
func (rr *RoleRepository) Create(role *model.Role) error {
	log.Infof("Creating role: %s", role.Id)
	return InsertOne(rr.collection, role)
}

// FindAll 查找所有角色
func (rr *RoleRepository) FindAll() ([]*model.Role, error) {
	log.Info("Finding all roles")
	results, err := FindMany[*model.Role](rr.collection, bson.M{})
	if err != nil {
		return nil, err
	}
	return results, nil
}
