package contollers

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gospotify.com/types"
	"gospotify.com/utils"
)

type User = types.User

func UsersController(root *gin.RouterGroup, db *mongo.Database) {
	usersRouter := root.Group("/users")

	usersRouter.GET("", func(ctx *gin.Context) {
		var users []User
		usersColl := db.Collection("users")
		cursor, usersErr := usersColl.Find(context.TODO(), gin.H{"is_artist": false})
		if usersErr != nil {
			panic(usersErr)
		}
		if usersErr := cursor.All(context.TODO(), &users); usersErr != nil {
			utils.HandleError(ctx, usersErr, "Failed to find users")
			panic(usersErr)
		}

		utils.JsonResponseOk(ctx, users)
	})

	usersRouter.GET("/artists", func(ctx *gin.Context) {
		var artists []User
		usersColl := db.Collection("users")
		cursor, usersErr := usersColl.Find(context.TODO(), gin.H{"is_artist": true})
		if usersErr != nil {
			panic(usersErr)
		}
		if usersErr := cursor.All(context.TODO(), &artists); usersErr != nil {
			ctx.JSON(500, gin.H{"error": usersErr})
			panic(usersErr)
		}
		utils.JsonResponseOk(ctx, artists)
	})
}
