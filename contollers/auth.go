package contollers

import (
	"context"
	"crypto/rand"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"gospotify.com/db"
	"gospotify.com/models"
	"gospotify.com/utils"
)

type RegUserForm = models.RegisterUserForm

func AuthController(r *gin.Engine, jwtMid *jwt.GinJWTMiddleware) {
	r.POST("/register", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		var userRegJson RegUserForm
		if err := c.ShouldBind(&userRegJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password is empty"})
			return
		}
		salt := make([]byte, 32)
		_, err := rand.Read(salt)
		if err != nil {
			panic(err)
		}
		passwordSalt, _ := utils.BytesToHex(salt)
		passwordHash, err := utils.HashSha256(password + passwordSalt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			panic(err)
		}
		user := models.RegisterUserDb{
			Username: username,
			Password: passwordHash,
			Salt:     passwordSalt,
		}
		if err := db.Db.Collection("users").FindOne(context.TODO(), bson.M{"username": user.Username}).Err(); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}

		cursor, err := db.Db.Collection("users").InsertOne(context.TODO(), user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create user: " + err.Error()})
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"id": cursor.InsertedID,
		})

	})
	r.POST("/login", jwtMid.LoginHandler)
	r.POST("/logout", jwtMid.LogoutHandler)
	r.NoRoute(jwtMid.MiddlewareFunc(), utils.HandleNoRoute())

	authorized := r.Group("/auth", jwtMid.MiddlewareFunc())
	TokenController(authorized, jwtMid)
	authorized.GET("/hello", HelloHandler)

	AdminController(authorized)
}

func TokenController(r *gin.RouterGroup, jwtMid *jwt.GinJWTMiddleware) {
	r.POST("/verify_token", func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		c.JSON(http.StatusOK, gin.H{
			"username": claims["id"],
		})
	})
	r.GET("/refresh_token", jwtMid.RefreshHandler)
}
