package admin

import (
	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gospotify.com/auth"
	cauth "gospotify.com/auth"
	"gospotify.com/contollers"
	"gospotify.com/db"
)

func AdminController(r *gin.RouterGroup) {
	jwtAdminMid, err := jwt.New(auth.InitJwtParams())
	if err != nil {
		log.Fatal("JWT error: ", err.Error())
	}
	jwtAdminMid.Authorizator = cauth.AdminAuthorizator()
	admin := r.Group("/admin", jwtAdminMid.MiddlewareFunc())
	contollers.UsersController(admin, db.Db)
}
