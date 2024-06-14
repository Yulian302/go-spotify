package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error, message string) {
	log.Println(err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}
