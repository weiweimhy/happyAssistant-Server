package repository

import (
	"happyAssistant/internal/initialize"
	"happyAssistant/internal/model"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// UserRepository 用户数据访问层
// 提供用户相关的数据库操作，如创建、查询、更新用户信息等
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository() *UserRepository {
	client := initialize.GetMongoClient()
	collection := client.Collection("users")
	return &UserRepository{
		collection: collection,
	}
}

// Create 创建新用户
func (ur *UserRepository) Create(user *model.User) error {
	log.Infof("Creating user: %s", user.Id)
	return InsertOne(ur.collection, user)
}

// FindByID 根据用户ID查找用户
func (ur *UserRepository) FindByID(userID string) (*model.User, error) {
	log.Infof("Finding user by ID: %s", userID)
	filter := bson.M{"_id": userID}
	result, err := FindOne[*model.User](ur.collection, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindByOpenID 根据微信OpenID查找用户
func (ur *UserRepository) FindByOpenID(openID string) (*model.User, error) {
	log.Infof("Finding user by OpenID: %s", openID)
	filter := bson.M{"open_id": openID}
	result, err := FindOne[*model.User](ur.collection, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Update 更新用户信息
func (ur *UserRepository) Update(user *model.User) error {
	log.Infof("Updating user: %s", user.Id)
	filter := bson.M{"_id": user.Id}
	return ReplaceOne(ur.collection, filter, user)
}

// Delete 删除用户
func (ur *UserRepository) Delete(userID string) error {
	log.Infof("Deleting user: %s", userID)
	filter := bson.M{"_id": userID}
	return DeleteOne(ur.collection, filter)
}

// FindAll 查找所有用户
func (ur *UserRepository) FindAll() ([]*model.User, error) {
	log.Info("Finding all users")
	results, err := FindMany[*model.User](ur.collection, bson.M{})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// FindByLabID 根据实验室ID查找用户
func (ur *UserRepository) FindByLabID(labID string) ([]*model.User, error) {
	log.Infof("Finding users by lab ID: %s", labID)
	filter := bson.M{"lib_ids": labID}
	results, err := FindMany[*model.User](ur.collection, filter)
	if err != nil {
		return nil, err
	}
	return results, nil
}
