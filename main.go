package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gofor-little/env"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gospotify.com/auth"
	"gospotify.com/contollers"
	"gospotify.com/types"
)

const ENV_FILE_PATH = ".env"

// importing types

type UserLogin = types.Login

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

	// root router
	router := gin.New()

	// jwt middleware
	authMiddleware, err := jwt.New(auth.InitJwtParams())
	if err != nil {
		log.Fatal("JWT error: ", err.Error())
	}
	router.Use(auth.HandlerMiddleWare(authMiddleware))
	auth.RegisterRoute(router, authMiddleware)

	// api routes group
	apiRouter := router.Group("/api/v1")

	apiRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// users routes
	contollers.UsersController(apiRouter, db)

	// songs routes
	contollers.SongsController(apiRouter, db)

	router.Run("localhost:8081")
}
