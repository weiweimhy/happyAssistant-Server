package repository

import (
	"context"
	"fmt"
	"happyAssistant/internal/config"
	"reflect"

	"google.golang.org/protobuf/proto"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// InsertOne 插入单个文档
// T 必须是 proto.Message 的指针类型
func InsertOne[T proto.Message](collection *mongo.Collection, data T) error {
	// 检查data是否为nil - 使用反射检查，因为泛型类型不能直接与nil比较
	if reflect.ValueOf(data).IsNil() {
		return fmt.Errorf("cannot insert nil document")
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg.MongoDB.OpTimeout)
	defer cancel()

	// 将泛型类型转换为interface{}，MongoDB驱动会处理序列化
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	return nil
}

// InsertMany 插入多个文档
// T 必须是 proto.Message 的指针类型
func InsertMany[T proto.Message](collection *mongo.Collection, data []T) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg.MongoDB.OpTimeout)
	defer cancel()

	// 检查data是否为空或包含nil元素
	if len(data) == 0 {
		return fmt.Errorf("cannot insert empty slice")
	}

	// 检查是否有nil元素
	for i, item := range data {
		if reflect.ValueOf(item).IsNil() {
			return fmt.Errorf("cannot insert nil document at index %d", i)
		}
	}

	// 直接使用data，因为T已经是proto.Message类型
	docs := make([]interface{}, len(data))
	for i, item := range data {
		docs[i] = item
	}
	_, err := collection.InsertMany(ctx, docs)
	return err
}

// DeleteOne 删除单个文档
func DeleteOne(collection *mongo.Collection, filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg.MongoDB.OpTimeout)
	defer cancel()
	_, err := collection.DeleteOne(ctx, filter)
	return err
}

// DeleteMany 删除多个文档
func DeleteMany(collection *mongo.Collection, filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg.MongoDB.OpTimeout)
	defer cancel()
	_, err := collection.DeleteMany(ctx, filter)
	return err
}

// FindOne 查找单个文档
// T 必须是 proto.Message 的指针类型
func FindOne[T proto.Message](collection *mongo.Collection, filter interface{}) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg.MongoDB.OpTimeout)
	defer cancel()

	// 创建T类型的零值
	var result T
	// 使用反射创建T的新实例
	result = proto.Clone(result).(T)

	err := collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		var zero T
		return zero, err
	}
	return result, nil
}

// FindMany 查找多个文档
// T 必须是 proto.Message 的指针类型
func FindMany[T proto.Message](collection *mongo.Collection, filter interface{}) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg.MongoDB.OpTimeout)
	defer cancel()

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Println("cursor close error:", err)
		}
	}(cur, ctx)

	var results []T
	for cur.Next(ctx) {
		// 创建T类型的新实例
		var elem T
		elem = proto.Clone(elem).(T)

		if err := cur.Decode(elem); err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// UpdateOne 更新单个文档
func UpdateOne(collection *mongo.Collection, filter interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg.MongoDB.OpTimeout)
	defer cancel()
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateMany 更新多个文档
func UpdateMany(collection *mongo.Collection, filter interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg.MongoDB.OpTimeout)
	defer cancel()
	_, err := collection.UpdateMany(ctx, filter, update)
	return err
}

// ReplaceOne 替换单个文档
// T 必须是 proto.Message 的指针类型
func ReplaceOne[T proto.Message](collection *mongo.Collection, filter interface{}, replacement T) error {
	// 检查replacement是否为nil
	if reflect.ValueOf(replacement).IsNil() {
		return fmt.Errorf("cannot replace with nil document")
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg.MongoDB.OpTimeout)
	defer cancel()
	_, err := collection.ReplaceOne(ctx, filter, replacement)
	return err
}

// Count 统计文档数量
func Count(collection *mongo.Collection, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Cfg.MongoDB.OpTimeout)
	defer cancel()
	return collection.CountDocuments(ctx, filter)
}
