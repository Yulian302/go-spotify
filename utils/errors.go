package utils

import (
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error, message string) {
	log.Println(err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}

// 404
func HandleNoRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}
