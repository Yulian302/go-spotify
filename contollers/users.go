package contollers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gospotify.com/models"
	"gospotify.com/utils"
)

type User = models.UserDb
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
