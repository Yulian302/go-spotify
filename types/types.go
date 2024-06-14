package types

import "time"

type User struct {
	Username   string    `bson:"username"`
	DateJoined time.Time `bson:"date_joined"`
	IsArtist   bool      `bson:"is_artist"`
}

type Song struct {
	Title     string   `bson:"title"`
	ArtistId  string   `bson:"artist_id"`
	AlbumId   string   `bson:"album_id"`
	Duration  float64  `bson:"duration"`
	SourceUrl string   `bson:"source_url"`
	Genres    []string `bson:"genres"`
	Year      int32    `bson:"year"`
	BitRate   int16    `bson:"bit_rate"`
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
