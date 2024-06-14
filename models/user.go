package models

import "time"

type User struct {
	Username   string    `bson:"username"`
	Password   string    `bson:"password"`
	DateJoined time.Time `bson:"date_joined"`
	IsArtist   bool      `bson:"is_artist"`
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
