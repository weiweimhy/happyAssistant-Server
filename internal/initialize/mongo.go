package initialize

import (
	"context"
	log "github.com/sirupsen/logrus"
	"happyAssistant/internal/config"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBClient struct {
	*mongo.Client
	dbName string
}

var (
	mongoClient *MongoDBClient
	once        sync.Once
)

func InitMongoDBClient(config config.MongoConfig) *MongoDBClient {
	once.Do(func() {
		opts := options.Client().
			SetConnectTimeout(config.Timeout).
			SetBSONOptions(&options.BSONOptions{
				UseJSONStructTags: true,
				NilSliceAsEmpty:   true,
				NilMapAsEmpty:     true,
				OmitEmpty:         true,
			}).
			ApplyURI(config.URI)

		if config.Username != "" && config.Password != "" {
			auth := options.Credential{
				Username:   config.Username,
				Password:   config.Password,
				AuthSource: config.Database,
			}
			opts.SetAuth(auth)
		}

		client, err := mongo.Connect(opts)
		if err != nil {
			log.Fatalln("Mongo connection error:", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
		defer cancel()
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatalln("Mongo ping error:", err)
		}

		mongoClient = &MongoDBClient{client, config.Database}

		log.Println("Mongo connection success")
	})
	return mongoClient
}

func GetMongoClient() *MongoDBClient {
	return mongoClient
}

func (c *MongoDBClient) Collection(collName string) *mongo.Collection {
	return c.Database(c.dbName).Collection(collName)
}
