package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JsonResponseOk(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}
