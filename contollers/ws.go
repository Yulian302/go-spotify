package contollers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func HandleWs(c *gin.Context) {
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
