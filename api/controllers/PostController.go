package controllers

import (
	entities "api/api/entities"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type postController struct {
	collection *mongo.Collection
}

func NewPostController(client *mongo.Client, dbName, collectionName string) *postController {
	collection := client.Database(dbName).Collection(collectionName)
	return &postController{collection}
}

func (p *postController) CreatePost(ctx *gin.Context) {
	var post entities.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post.ID = uuid.New()
	post.PostAt = time.Now().Format("2006-01-02 15:04")

	_, err := p.collection.InsertOne(context.TODO(), post)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, post)
}

func (p *postController) FindAll(ctx *gin.Context) {
	cursor, err := p.collection.Find(context.TODO(), bson.M{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var posts []entities.Post
	for cursor.Next(context.TODO()) {
		var post entities.Post
		err := cursor.Decode(&post)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		posts = append(posts, post)
	}

	ctx.JSON(http.StatusOK, posts)
}

func (p *postController) DeletePost(ctx *gin.Context) {
	postID := ctx.Param("id")
	fmt.Println(postID)

	if uuid.Parse(postID) == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	filter := bson.M{"id": postID}

	result, err := p.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Postagem não encontrada"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Postagem excluída com sucesso"})

}
