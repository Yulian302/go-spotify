package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gospotify.com/auth"
	"gospotify.com/contollers"
	env "gospotify.com/env"
)

func main() {
	env.LoadEnv()
	// mongodb config and connection
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@%s.hfa617f.mongodb.net/?retryWrites=true&w=majority&appName=%s", env.ClusterName, env.UserPassword, strings.ToLower(env.ClusterName), env.ClusterName)).SetServerAPIOptions(serverAPI)

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
	db := client.Database(env.DbName)

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
