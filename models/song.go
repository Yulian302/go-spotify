package models

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
