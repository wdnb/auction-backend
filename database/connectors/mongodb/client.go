package mdb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

var client *mongo.Client

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
}

func GetClient(c *Config) *mongo.Client {
	if client == nil {
		client = initClient(c)
	}
	return client
}

func initClient(c *Config) *mongo.Client {
	// 构建连接字符串
	uri := "mongodb://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port
	// 设置客户端选项
	opts := options.Client().ApplyURI(uri)

	// 连接重试逻辑
	client, err := connectWithRetry(context.Background(), opts)
	if err != nil {
		panic(err)
	}

	return client
}

func connectWithRetry(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
	var maxRetries = 3               // 最大重试次数
	retryInterval := 5 * time.Second // 重试间隔

	for i := 0; i < maxRetries; i++ {
		client, err := mongo.Connect(ctx, opts...)
		if err == nil {
			return client, nil
		}

		zap.L().Error("Failed to connect to MongoDB",
			zap.Int("attempt", i+1),
			zap.Int("max_retries", maxRetries),
			zap.Error(err),
		)

		if i == maxRetries-1 {
			return nil, fmt.Errorf("failed to connect after %d retries", maxRetries)
		}

		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("should not reach here")
}

//func test(app *utils.Components) {
//	// 获取指定的数据库和集合
//	database := app.MongoDB.Database(viper.GetString("mongodb.dbname"))
//	collection := database.Collection("test")
//	// 插入一个文档
//	document := bson.M{"name": "John Doe", "age": 30}
//	err := InsertDocument(collection, document)
//	if err != nil {
//		fmt.Println("Failed to insert document:", err)
//		return
//	}
//
//	// 更新一个文档
//	filter := bson.M{"name": "John Doe"}
//	update := bson.M{"$set": bson.M{"age": 31}}
//	err = UpdateDocument(collection, filter, update)
//	if err != nil {
//		fmt.Println("Failed to update document:", err)
//		return
//	}
//
//	// 删除一个文档
//	//filter = bson.M{"name": "John Doe"}
//	//err = DeleteDocument(collection, filter)
//	//if err != nil {
//	//	fmt.Println("Failed to delete document:", err)
//	//	return
//	//}
//
//	// 读取集合中的所有文档
//	results, err := ReadAllDocuments(collection)
//	if err != nil {
//		fmt.Println("Failed to read documents:", err)
//		return
//	}
//	fmt.Println(results)
//
//}

// 读取集合中的所有文档
func ReadAllDocuments(collection *mongo.Collection) ([]bson.M, error) {
	var results []bson.M
	// 执行查询操作
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return results, err
	}
	// 将结果解码为 bson.M 类型的文档
	if err := cursor.All(context.Background(), &results); err != nil {
		return results, err
	}
	return results, nil
}

// 向集合中插入一个文档
func InsertDocument(collection *mongo.Collection, document interface{}) error {
	// 执行插入操作
	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		return err
	}
	return nil
}

// 更新集合中的一个文档
func UpdateDocument(collection *mongo.Collection, filter interface{}, update interface{}) error {
	// 执行更新操作
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

// 删除集合中的一个文档
func DeleteDocument(collection *mongo.Collection, filter interface{}) error {
	// 执行删除操作
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
