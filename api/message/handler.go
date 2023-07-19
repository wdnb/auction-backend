package message

import (
	"auction-website/conf"
	"auction-website/message/im"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type Handler struct {
	Service *Service
}

func NewHandler(c *conf.Config) *Handler {
	return &Handler{
		Service: NewService(c),
	}
}

func (h *Handler) WsHandler(c *gin.Context) {
	uid := c.Query("uid")      // 自己的id
	toUid := c.Query("to_uid") // 对方的id
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
			return true
		}}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	// 创建一个用户实例
	client := NewClient(uid, toUid, conn, h.Service.repo.rdb, h.Service.repo.mdb)
	// 用户注册到用户管理上
	im.Manager.Register <- client
	go client.Read()
	go client.Write()
}

func NewClient(uid, toUid string, conn *websocket.Conn, redisClient *redis.Client, mongoClient *mongo.Client) *im.Client {
	return &im.Client{
		ID:            createId(uid, toUid),
		SendID:        createId(toUid, uid),
		Socket:        conn,
		Redis:         redisClient,
		MongoDBClient: mongoClient,
		Send:          make(chan []byte),
	}
}

func createId(uid, toUid string) string {
	return uid + "->" + toUid
}
