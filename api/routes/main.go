package routes

import (
	controllers "api/api/controllers"
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

	postController := controllers.NewPostController(client, "Blinkr", "posts")
	userController := controllers.NewUserController(client, "Blinkr", "users")

	v1 := router.Group("/v1")
	{

		// POSTS
		v1.GET("/posts", postController.FindAll)
		v1.POST("/posts", postController.CreatePost)
		v1.DELETE("/posts/:id", postController.DeletePost)

		// USERS

		v1.GET("/user", userController.FindAllUsers)
		v1.POST("/user", userController.CreateUser)
		v1.DELETE("/user/:id", userController.DeleteUser)
		v1.PATCH("/user/:id", userController.UpdateUser)

	}

	return v1
}
