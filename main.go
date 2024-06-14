package main

import (
	"context"
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gospotify.com/auth"
	"gospotify.com/contollers"
	"gospotify.com/db"
	env "gospotify.com/env"
)

func main() {
	env.LoadEnv()

	// init db connection
	client, err := db.DbClient()
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}
	// choosing database
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
