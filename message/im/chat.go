package im

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strconv"
	"time"
)

const month = 60 * 60 * 24 * 30 // 按照30天算一个月

const (
	SUCCESS               = 200
	UpdatePasswordSuccess = 201
	NotExistInentifier    = 202
	ERROR                 = 500
	InvalidParams         = 400
	ErrorDatabase         = 40001

	WebsocketSuccessMessage = 50001
	WebsocketSuccess        = 50002
	WebsocketEnd            = 50003
	WebsocketOnlineReply    = 50004
	WebsocketOfflineReply   = 50005
	WebsocketLimit          = 50006
)

// 发送消息的类型
type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// 回复的消息
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

// 用户类
type Client struct {
	ID            string
	SendID        string
	Socket        *websocket.Conn
	Redis         *redis.Client
	MongoDBClient *mongo.Client
	Send          chan []byte
}

// 广播类，包括广播内容和源用户
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// 用户管理
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

// Message 信息转JSON (包括：发送者、接收者、内容)
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}

func (c *Client) Read() {
	defer func() { // 避免忘记关闭，所以要加上close
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		sendMsg := new(SendMsg)
		//_,msg,_:=c.Socket.ReadMessage()
		err := c.Socket.ReadJSON(&sendMsg) // 读取json格式，如果不是json格式，会报错
		if err != nil {
			log.Println("数据格式不正确", err)
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}
		//fmt.Println(sendMsg)
		var ctx = context.Background()
		if sendMsg.Type == 1 {
			r1, _ := c.Redis.Get(ctx, c.ID).Result()
			r2, _ := c.Redis.Get(ctx, c.SendID).Result()
			if r1 >= "3" && r2 == "" { // 限制单聊
				replyMsg := ReplyMsg{
					Code:    WebsocketLimit,
					Content: "达到限制",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				_, _ = c.Redis.Expire(ctx, c.ID, time.Hour*24*30).Result() // 防止重复骚扰，未建立连接刷新过期时间一个月
				continue
			} else {
				c.Redis.Incr(ctx, c.ID)
				_, _ = c.Redis.Expire(ctx, c.ID, time.Hour*24*30*3).Result() // 防止过快“分手”，建立连接三个月过期
			}
			log.Println(c.ID, "发送消息", sendMsg.Content)
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(sendMsg.Content),
			}
		} else if sendMsg.Type == 2 { //拉取历史消息
			timeT, err := strconv.Atoi(sendMsg.Content) // 传送来时间
			if err != nil {
				timeT = 999999999
			}
			results, _ := c.FindMany(viper.GetString("mongodb.dbname"), c.SendID, c.ID, int64(timeT), 10)
			if len(results) > 10 {
				results = results[:10]
			} else if len(results) == 0 {
				replyMsg := ReplyMsg{
					Code:    WebsocketEnd,
					Content: "到底了",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: fmt.Sprintf("%s", result.Msg),
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		} else if sendMsg.Type == 3 {
			results, err := c.FirsFindtMsg(viper.GetString("mongodb.dbname"), c.SendID, c.ID)
			if err != nil {
				log.Println(err)
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: fmt.Sprintf("%s", result.Msg),
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}

func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Println(c.ID, "接受消息:", string(message))
			replyMsg := ReplyMsg{
				Code:    WebsocketSuccessMessage,
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
