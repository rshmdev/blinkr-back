package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WebSocketController struct {
	collection *mongo.Collection
	conns      []*websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Sender    string
	Receiver  string
	Message   string
	Timestamp time.Time
}

var Users []websocket.Conn

func NewWebSocketController(client *mongo.Client, dbName, collectionName string) *WebSocketController {
	collection := client.Database(dbName).Collection(collectionName)
	return &WebSocketController{collection: collection}
}

func (wsc *WebSocketController) WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	wsc.conns = append(wsc.conns, conn)

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var messageObject Message
		json.Unmarshal(msg, &messageObject)

		insertResult, err := wsc.collection.InsertOne(context.TODO(), bson.D{
			{Key: "sender", Value: messageObject.Sender},
			{Key: "receiver", Value: messageObject.Receiver},
			{Key: "message", Value: messageObject.Message},
			{Key: "timestamp", Value: time.Now()},
		})

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Inserted a single document: ", insertResult.InsertedID)

		for _, conn := range wsc.conns {
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	}
}

func (m *WebSocketController) GetMessages(c *gin.Context) {
	// Recupere o ID do usuário dos parâmetros da URL
	userId := c.Query("userId")

	// Crie um filtro para recuperar as mensagens corretas
	var filter bson.M
	if userId != "" {
		// Se um ID de usuário foi fornecido, recupere apenas as mensagens enviadas por esse usuário
		filter = bson.M{"sender": userId}
	} else {
		// Se nenhum ID de usuário foi fornecido, recupere todas as mensagens
		filter = bson.M{}
	}

	// Recupere todas as mensagens que correspondem ao filtro
	cursor, err := m.collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	var messages []bson.M
	if err = cursor.All(context.TODO(), &messages); err != nil {
		log.Fatal(err)
	}

	// Retorne as mensagens como resposta JSON
	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
