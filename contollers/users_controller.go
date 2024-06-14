package contollers

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gospotify.com/models"
	"gospotify.com/utils"
)

type User = models.User

func UsersController(root *gin.RouterGroup, db *mongo.Database) {
	usersRouter := root.Group("/users")

	usersRouter.GET("", func(c *gin.Context) {
		var users []User
		usersColl := db.Collection("users")
		cursor, usersErr := usersColl.Find(context.TODO(), gin.H{"is_artist": false})
		if usersErr != nil {
			panic(usersErr)
		}
		if usersErr := cursor.All(context.TODO(), &users); usersErr != nil {
			utils.HandleError(c, usersErr, "Failed to find users")
			panic(usersErr)
		}

		utils.JsonResponseOk(c, users)
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
