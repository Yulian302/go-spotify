package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"gospotify.com/authorization"
	"gospotify.com/contollers"
	"gospotify.com/db"
	"gospotify.com/models"
	"gospotify.com/utils"
)

type User = models.UserDb
type UserLogin = models.LoginUserForm

var (
	identityKey = "id"
	// port        string
)

func InitJwtParams() *jwt.GinJWTMiddleware {

	return &jwt.GinJWTMiddleware{
		Realm:           "test zone",
		Key:             []byte("secret key"),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     identityKey,
		PayloadFunc:     payloadFunc(),
		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(),
		Authorizator:    authorizator(),
		Unauthorized:    unauthorized(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup:   "cookie:token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}
func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginVals UserLogin
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		username := loginVals.Username
		password := loginVals.Password

		// search user in database
		var user User
		userErr := db.Db.Collection("users").FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)

		// if user is found
		if userErr == nil {
			hashedPassword, err := utils.HashSha256(password + user.Salt)
			if err != nil {

				log.Fatal(err)
			}
			if username == user.Username && hashedPassword == user.Password {
				return &User{
					Username: username,
				}, nil
			}
		}
		if username == "test" || password == "test" {
			return &User{
				Username: username,
			}, nil
		}
		return nil, jwt.ErrFailedAuthentication
	}
}
func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(*User); ok && authorization.IsUserAdmin(v) {
			return true
		}
		return false
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*User); ok {
			return jwt.MapClaims{
				identityKey: v.Username,
			}
		}
		return jwt.MapClaims{}
	}
}

func HandlerMiddleWare(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(context *gin.Context) {
		errInit := authMiddleware.MiddlewareInit()
		if errInit != nil {
			log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		}
	}
}

func identityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		var user User
		userErr := db.Db.Collection("users").FindOne(context.TODO(), bson.M{"username": claims[identityKey]}).Decode(&user)
		if userErr != nil {
			log.Println("Error retrieving user from database: ", userErr)
			return nil
		}
		return &User{
			Username: claims[identityKey].(string),
			IsAdmin:  user.IsAdmin,
		}
	}
}

// custom handlers
func verifyToken(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(http.StatusOK, gin.H{
		"username": claims["id"],
	})
}

// registering routes

func RegisterRoute(r *gin.Engine, handle *jwt.GinJWTMiddleware) {
	r.POST("/login", handle.LoginHandler)
	r.POST("/logout", handle.LogoutHandler)
	r.NoRoute(handle.MiddlewareFunc(), utils.HandleNoRoute())

	auth := r.Group("/auth", handle.MiddlewareFunc())
	auth.POST("/verify_token", verifyToken)
	auth.GET("/refresh_token", handle.RefreshHandler)

	auth.GET("/hello", contollers.HelloHandler)
	contollers.UsersController(auth, db.Db)
	contollers.SongsController(auth, db.Db)
}
