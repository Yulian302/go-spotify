package contollers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gospotify.com/models"
	"gospotify.com/utils"
)

type Song = models.Song

func SongsController(root *gin.RouterGroup, db *mongo.Database) {
	songsRouter := root.Group("/songs")
	songsRouter.GET("", func(c *gin.Context) {
		var songs []Song
		songsColl := db.Collection("songs")
		var queryCondition *gin.H = &gin.H{}
		if title := c.Query("title"); title != "" {
			queryCondition = &gin.H{"title": title}
		}
		cursor, songsErr := songsColl.Find(context.TODO(), queryCondition)
		if songsErr != nil {
			panic(songsErr)
		}
		if songsErr := cursor.All(context.TODO(), &songs); songsErr != nil {
			c.JSON(500, gin.H{"error": songsErr})
			panic(songsErr)
		}

		// if one element, flatten
		if len(songs) == 1 {
			utils.JsonResponseOk(c, songs[0])
			return
		}
		utils.JsonResponseOk(c, songs)
	})
	songsRouter.GET("/:id", func(c *gin.Context) {
		songId, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song id"})
			return
		}
		var song Song
		songError := db.Collection("songs").FindOne(context.TODO(), bson.M{
			"_id": songId,
		}).Decode(&song)
		if songError != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
			return
		}
		utils.JsonResponseOk(c, song)
	})
}
