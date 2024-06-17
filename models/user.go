package models

import "time"

type UserDb struct {
	Username   string    `bson:"username"`
	Password   string    `bson:"password"`
	Salt       string    `bson:"salt"`
	DateJoined time.Time `bson:"date_joined"`
	IsArtist   bool      `bson:"is_artist"`
}

type LoginUserForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RegisterUserForm struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type RegisterUserDb struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	Salt     string `bson:"salt"`
}
