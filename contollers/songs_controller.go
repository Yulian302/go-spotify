package contollers

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gospotify.com/models"
	"gospotify.com/utils"
)

type Song = models.Song

func SongsController(root *gin.RouterGroup, db *mongo.Database) {
	songsRouter := root.Group("/songs")
	songsRouter.GET("", func(ctx *gin.Context) {
		var songs []Song
		songsColl := db.Collection("songs")
		cursor, songsErr := songsColl.Find(context.TODO(), gin.H{})
		if songsErr != nil {
			panic(songsErr)
		}
		if songsErr := cursor.All(context.TODO(), &songs); songsErr != nil {
			ctx.JSON(500, gin.H{"error": songsErr})
			panic(songsErr)
		}
		utils.JsonResponseOk(ctx, songs)
	})
}
