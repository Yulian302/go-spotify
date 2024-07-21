package contollers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
func BasicController(r *gin.Engine){
	
}
func HandlePong(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

func HelloHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello",
	})
}
