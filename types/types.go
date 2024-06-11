package types

import "time"

type User struct {
	Username   string    `bson:"username"`
	DateJoined time.Time `bson:"date_joined"`
	IsArtist   bool      `bson:"is_artist"`
}
