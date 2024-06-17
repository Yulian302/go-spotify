package main

import (
	"context"

	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gospotify.com/auth"
	"gospotify.com/contollers"
	"gospotify.com/db"
	env "gospotify.com/env"
)

func handleWs(c *gin.Context) {
	upgrader := websocket.Upgrader{}
	// TODO only for dev
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	conn.WriteMessage(websocket.TextMessage, []byte("Hello this is ws"))
	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("Received close message from client")
				break
			} else {
				log.Println("Error reading message:", err)
				break
			}
		}
		if string(message) == "close" {
			conn.WriteMessage(websocket.CloseMessage, []byte{})
			break
		}
		log.Printf("Received message: %s with type: %d\n", message, messageType)
	}
}
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

	router.GET("/ws", handleWs)

	// jwt middleware
	authMiddleware, err := jwt.New(auth.InitJwtParams())
	if err != nil {
		log.Fatal("JWT error: ", err.Error())
	}

	// public routes
	// api routes group
	apiRouterPublic := router.Group("/api/v1")
	// songs routes
	contollers.SongsController(apiRouterPublic, db.Db)

	apiRouterPublic.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	apiRouterPublic.POST("/register", contollers.RegisterHandler)

	// private routes
	apiRouterPrivate := router.Group("/api/v1")
	apiRouterPrivate.Use(auth.HandlerMiddleWare(authMiddleware))
	auth.RegisterRoute(router, authMiddleware)

	// users routes
	contollers.UsersController(apiRouterPublic, db.Db)

	router.Run("localhost:8081")
}
