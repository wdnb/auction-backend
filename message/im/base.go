package im

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// 读取集合中的所有文档
func ReadAllDocument(collection *mongo.Collection, filter bson.M, order bson.D, limit int64) ([]Trainer, error) {
	var results []Trainer

	// 执行查询操作
	cursor, err := collection.Find(
		context.Background(),
		filter,
		options.Find().SetSort(order),
		options.Find().SetLimit(limit),
	)
	if err != nil {
		return results, err
	}
	// 将结果解码为 bson.M 类型的文档
	if err := cursor.All(context.Background(), &results); err != nil {
		return results, err
	}
	//fmt.Println(results)
	return results, nil
}

func InsertMsg(c *mongo.Client, database string, id string, content string, read uint, expire int64) (err error) {
	collection := c.Database(database).Collection(id)
	comment := Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      read,
	}
	_, err = collection.InsertOne(context.TODO(), comment)
	return
}

func UpdateAllDocument(collection *mongo.Collection, filter bson.M, update bson.M) error {
	// 执行更新操作
	_, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}
	//fmt.Println(result.ModifiedCount)
	return nil
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
