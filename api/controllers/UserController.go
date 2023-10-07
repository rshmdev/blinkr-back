package controllers

import (
	"api/api/entities"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userController struct {
	collection *mongo.Collection
}

func NewUserController(client *mongo.Client, dbName, collectionName string) *userController {
	collection := client.Database(dbName).Collection(collectionName)
	return &userController{collection}
}

func (p *userController) checkIfFieldExists(field, value string) (bool, error) {
	filter := bson.M{field: value}
	count, err := p.collection.CountDocuments(context.TODO(), filter)
	return count > 0, err
}

func (p *userController) CreateUser(ctx *gin.Context) {
	var user entities.User

	// if err := ctx.ShouldBindJSON(&user); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	if emailExist, err := p.checkIfFieldExists("email", user.Email); err != nil || emailExist {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email já está em uso"})
		return
	}

	if phoneExist, err := p.checkIfFieldExists("phone", user.Phone); err != nil || phoneExist {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Número de telefone já está em uso"})
		return
	}

	if usernameExist, err := p.checkIfFieldExists("username", user.Username); err != nil || usernameExist {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Nome de usuário já está em uso"})
		return
	}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = uuid.New()

	_, err := p.collection.InsertOne(context.TODO(), user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (p *userController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	var updatedUser entities.User
	if err := ctx.ShouldBindJSON(&updatedUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"id": userID}
	existingUser := p.collection.FindOne(context.TODO(), filter)
	if existingUser.Err() != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	var userToUpdate entities.User
	existingUser.Decode(&userToUpdate)

	if updatedUser.Email != "" {
		userToUpdate.Email = updatedUser.Email
	}

	if updatedUser.Username != "" {
		userToUpdate.Username = updatedUser.Username
	}

	if updatedUser.FirstName != "" {
		userToUpdate.FirstName = updatedUser.FirstName
	}

	if updatedUser.LastName != "" {
		userToUpdate.LastName = updatedUser.LastName
	}

	if updatedUser.Phone != "" {
		userToUpdate.Phone = updatedUser.Phone
	}

	if updatedUser.Birthday != "" {
		userToUpdate.Birthday = updatedUser.Birthday
	}

	if updatedUser.Avatar != "" {
		userToUpdate.Avatar = updatedUser.Avatar
	}

	update := bson.M{
		"$set": userToUpdate,
	}

	_, err := p.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Usuário atualizado com sucesso"})
}

func (p *userController) FindAllUsers(ctx *gin.Context) {
	cursor, err := p.collection.Find(context.TODO(), bson.M{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var users []entities.User
	for cursor.Next(context.TODO()) {
		var user entities.User
		err := cursor.Decode(&user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	ctx.JSON(http.StatusOK, users)
}

func (p *userController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")

	if uuid.Parse(userID) == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	filter := bson.M{"id": userID}

	result, err := p.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuario não encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Usuario excluído com sucesso"})

}
