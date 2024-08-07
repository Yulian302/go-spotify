package main

import (
	"context"

	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
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

	// root router
	router := gin.New()
	// CORS
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"*"},
		AllowAllOrigins:  true,
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	router.GET("/ping", contollers.HandlePong)
	router.GET("/ws", contollers.HandleWs)

	// jwt middleware
	authMiddleware, err := jwt.New(auth.InitJwtParams())
	if err != nil {
		log.Fatal("JWT error: ", err.Error())
	}

	contollers.AuthController(router, authMiddleware)

	// public routes
	apiRouterPublic := router.Group("/api/v1")

	// songs
	contollers.SongsController(apiRouterPublic, db.Db)

	// private routes
	apiRouterPrivate := router.Group("/api/v1")
	apiRouterPrivate.Use(auth.HandlerMiddleWare(authMiddleware))

	// users routes
	// contollers.UsersController(apiRouterPublic, db.Db)
	router.Run("localhost:8081")
}
