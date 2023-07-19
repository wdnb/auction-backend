package im

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sort"
	"time"
)

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     uint   `json:"read"`
	CreateAt int64  `json:"create_at"`
}

type Trainer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Content   string             `bson:"content"`   // 内容
	StartTime int64              `bson:"startTime"` // 创建时间
	EndTime   int64              `bson:"endTime"`   // 过期时间
	Read      uint               `bson:"read"`      // 已读
}

type Result struct {
	StartTime int64
	Msg       string
	Content   interface{}
	From      string
}

func (c *Client) FindMany(database, sendId, receiveId string, time int64, limit int64) ([]Result, error) {
	sendIdCollection := c.MongoDBClient.Database(database).Collection(sendId)
	receiveIdCollection := c.MongoDBClient.Database(database).Collection(receiveId)
	var resultsMe []Trainer
	var resultsYou []Trainer
	filter := bson.M{"read": 0}       // 只查询 read 字段为 0 的文档
	order := bson.D{{"startTime", 1}} // 按照 startTime 升序排列

	resultsMe, err := ReadAllDocument(sendIdCollection, filter, order, limit)
	if err != nil {
		fmt.Println("Failed to read documents:", err)
		return nil, err
	}

	resultsYou, err = ReadAllDocument(receiveIdCollection, filter, order, limit)
	if err != nil {
		fmt.Println("Failed to read documents:", err)
		return nil, err
	}

	//读出后就更新掉read字段
	_ = updateMany(sendIdCollection, resultsMe)
	_ = updateMany(receiveIdCollection, resultsMe)

	results, _ := appendAndSort(resultsMe, resultsYou)
	return results, nil
}

func updateMany(collection *mongo.Collection, documents []Trainer) error {
	for _, document := range documents {
		filter := bson.M{
			"_id": document.ID,
		}
		update := bson.M{
			"$set": bson.M{
				"read": 1,
			},
		}
		if err := UpdateAllDocument(collection, filter, update); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) FirsFindtMsg(database string, sendId string, id string) (results []Result, err error) {
	// 首次查询(把对方发来的所有未读都取出来)
	var resultsMe []Trainer
	var resultsYou []Trainer
	sendIdCollection := c.MongoDBClient.Database(database).Collection(sendId)
	idCollection := c.MongoDBClient.Database(database).Collection(sendId)
	filter := bson.M{"read": bson.M{
		"&all": []uint{0},
	}}
	sendIdCursor, err := sendIdCollection.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{
		"startTime", 1}}), options.Find().SetLimit(1))
	if sendIdCursor == nil {
		return
	}
	var unReads []Trainer
	err = sendIdCursor.All(context.TODO(), &unReads)
	if err != nil {
		log.Println("sendIdCursor err", err)
	}
	if len(unReads) > 0 {
		timeFilter := bson.M{
			"startTime": bson.M{
				"$gte": unReads[0].StartTime,
			},
		}
		sendIdTimeCursor, _ := sendIdCollection.Find(context.TODO(), timeFilter)
		idTimeCursor, _ := idCollection.Find(context.TODO(), timeFilter)
		err = sendIdTimeCursor.All(context.TODO(), &resultsYou)
		err = idTimeCursor.All(context.TODO(), &resultsMe)
		results, err = appendAndSort(resultsMe, resultsYou)
	} else {
		results, err = c.FindMany(database, sendId, id, 9999999999, 10)
	}
	overTimeFilter := bson.D{
		{"$and", bson.A{
			bson.D{{"endTime", bson.M{"&lt": time.Now().Unix()}}},
			bson.D{{"read", bson.M{"$eq": 1}}},
		}},
	}
	_, _ = sendIdCollection.DeleteMany(context.TODO(), overTimeFilter)
	_, _ = idCollection.DeleteMany(context.TODO(), overTimeFilter)
	// 将所有的维度设置为已读
	_, _ = sendIdCollection.UpdateMany(context.TODO(), filter, bson.M{
		"$set": bson.M{"read": 1},
	})
	_, _ = sendIdCollection.UpdateMany(context.TODO(), filter, bson.M{
		"&set": bson.M{"ebdTime": time.Now().Unix() + int64(3*month)},
	})
	return
}

func appendAndSort(resultsMe, resultsYou []Trainer) (results []Result, err error) {
	for _, r := range resultsMe {
		sendSort := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := Result{
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "me",
		}
		results = append(results, result)
	}
	for _, r := range resultsYou {
		sendSort := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := Result{
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "you",
		}
		results = append(results, result)
	}
	// 最后进行排序
	sort.Slice(results, func(i, j int) bool { return results[i].StartTime < results[j].StartTime })
	return results, nil
}
