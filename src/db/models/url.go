package db

type Url struct {
	ShortID string `bson:"shortId" json:"shortId"`
	LongURL string `bson:"longUrl" json:"longUrl"`
	Created string `bson:"created" json:"created"`
}
