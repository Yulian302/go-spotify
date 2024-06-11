package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofor-little/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gospotify.com/types"
)

const ENV_FILE_PATH = ".env"

// importing types
type User = types.User

func main() {

	// environment variables
	if err := env.Load(ENV_FILE_PATH); err != nil {
		panic(err)
	}
	clusterName := env.Get("CLUSTER_NAME", "")
	userPassword := env.Get("CLUSTER_USER_PASSWD", "")
	dbName := env.Get("DB_NAME", "")

	// mongodb config and connection
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@%s.hfa617f.mongodb.net/?retryWrites=true&w=majority&appName=%s", clusterName, userPassword, strings.ToLower(clusterName), clusterName)).SetServerAPIOptions(serverAPI)

	// create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	// choosing db
	db := client.Database(dbName)

	router := gin.New()
	apiRouter := router.Group("/api/v1")
	apiRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	usersRouter := apiRouter.Group("/users")
	usersRouter.GET("/", func(ctx *gin.Context) {
		var users []User
		usersColl := db.Collection("users")
		cursor, usersErr := usersColl.Find(context.TODO(), gin.H{})
		if usersErr != nil {
			panic(err)
		}
		if usersErr := cursor.All(context.TODO(), &users); usersErr != nil {
			ctx.JSON(500, gin.H{"error": err})
			panic(usersErr)
		}
		ctx.JSON(200, gin.H{"data": users})
	})
	router.Run("localhost:8081")
}
