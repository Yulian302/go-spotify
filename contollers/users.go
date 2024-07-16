package contollers

import (
	"context"
	"crypto/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gospotify.com/db"
	"gospotify.com/models"
	"gospotify.com/utils"
)

type User = models.UserDb
type RegUserForm = models.RegisterUserForm
type RegUserDb = models.RegisterUserDb

func UsersController(root *gin.RouterGroup, db *mongo.Database) {
	usersRouter := root.Group("/users")

	usersRouter.GET("", func(c *gin.Context) {
		var users []User
		usersColl := db.Collection("users")
		var usersCondition *gin.H = &gin.H{"is_artist": false}
		if name := c.Query("username"); name != "" {
			(*usersCondition)["username"] = name
		}
		cursor, usersErr := usersColl.Find(context.TODO(), usersCondition)
		if usersErr != nil {
			panic(usersErr)
		}
		if usersErr := cursor.All(context.TODO(), &users); usersErr != nil {
			utils.HandleError(c, usersErr, "Failed to find users")
			panic(usersErr)
		}
		// if one element, flatten
		if len(users) == 1 {
			utils.JsonResponseOk(c, users[0])
			return
		}
		utils.JsonResponseOk(c, users)
	})

	usersRouter.GET("/:id", func(c *gin.Context) {
		usersColl := db.Collection("users")
		userId, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}
		var user *User
		userError := usersColl.FindOne(context.TODO(), gin.H{"_id": userId}).Decode(&user)
		if userError != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": userError})
		}

		utils.JsonResponseOk(c, user)
	})

	usersRouter.GET("/artists", func(c *gin.Context) {
		var artists []User
		usersColl := db.Collection("users")
		cursor, usersErr := usersColl.Find(context.TODO(), gin.H{"is_artist": true})
		if usersErr != nil {
			panic(usersErr)
		}
		if usersErr := cursor.All(context.TODO(), &artists); usersErr != nil {
			c.JSON(500, gin.H{"error": usersErr})
			panic(usersErr)
		}
		utils.JsonResponseOk(c, artists)
	})

}

func RegisterHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	var userRegJson RegUserForm
	if err := c.ShouldBind(&userRegJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password is empty"})
		return
	}
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}
	passwordSalt, _ := utils.BytesToHex(salt)
	passwordHash, err := utils.HashSha256(password + passwordSalt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		panic(err)
	}
	user := models.RegisterUserDb{
		Username: username,
		Password: passwordHash,
		Salt:     passwordSalt,
	}
	if err := db.Db.Collection("users").FindOne(context.TODO(), bson.M{"username": user.Username}).Err(); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	cursor, err := db.Db.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create user: " + err.Error()})
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"id": cursor.InsertedID,
	})
}
