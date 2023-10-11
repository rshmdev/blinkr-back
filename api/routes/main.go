package routes

import (
	controllers "api/api/controllers"
	"api/api/websocket"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AppRoutes(router *gin.Engine) *gin.RouterGroup {

	clientOptions := options.Client().ApplyURI("mongodb+srv://rianmoraes:rianmoraes@blinkr.vjzgh6b.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	userController := controllers.NewUserController(client, "Blinkr", "users")
	postController := controllers.NewPostController(client, "Blinkr", "posts", userController)
	wsController := websocket.NewWebSocketController(client, "Blinkr", "messages")

	v1 := router.Group("/v1")
	{

		// POSTS
		v1.GET("/posts", postController.FindAll)
		v1.POST("/posts", postController.CreatePost)
		v1.DELETE("/posts/:id", postController.DeletePost)

		// USERS

		v1.GET("/users", userController.FindAllUsers)
		v1.GET(("users/:id"), userController.GetUserById)
		v1.DELETE("/users/:id", userController.DeleteUser)
		v1.PATCH("/users/:id", userController.UpdateUser)

		// Auth

		v1.POST("/users/login", userController.Login)
		v1.POST("/users/register", userController.CreateUser)

		v1.GET("/ws", wsController.WebSocketHandler)
		v1.GET("/messages", wsController.GetMessages)

	}

	return v1
}
