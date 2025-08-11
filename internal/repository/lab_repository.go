package repository

import (
	"happyAssistant/internal/initialize"
	"happyAssistant/internal/model"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// LabRepository 实验室数据访问层
type LabRepository struct {
	collection *mongo.Collection
}

// NewLabRepository 创建实验室仓库实例
func NewLabRepository() *LabRepository {
	client := initialize.GetMongoClient()
	collection := client.Collection("labs")
	return &LabRepository{
		collection: collection,
	}
}

// GetDefaultLab 获取默认实验室
func (lr *LabRepository) GetDefaultLab() (*model.Lab, error) {
	log.Info("Getting default lab")
	filter := bson.M{"is_default": true}
	result, err := FindOne[*model.Lab](lr.collection, filter)
	if err != nil {
		// 如果没有默认实验室，返回第一个实验室
		filter = bson.M{}
		result, err = FindOne[*model.Lab](lr.collection, filter)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// FindByID 根据ID查找实验室
func (lr *LabRepository) FindByID(labID string) (*model.Lab, error) {
	log.Infof("Finding lab by ID: %s", labID)
	filter := bson.M{"_id": labID}
	result, err := FindOne[*model.Lab](lr.collection, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Create 创建实验室
func (lr *LabRepository) Create(lab *model.Lab) error {
	log.Infof("Creating lab: %s", lab.Id)
	return InsertOne(lr.collection, lab)
}

// FindAll 查找所有实验室
func (lr *LabRepository) FindAll() ([]*model.Lab, error) {
	log.Info("Finding all labs")
	results, err := FindMany[*model.Lab](lr.collection, bson.M{})
	if err != nil {
		return nil, err
	}
	return results, nil
}
