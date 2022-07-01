package db

import "time"

type User struct {
	//TODO ID       bson.ObjectId `bson:"_id" json:"id" example:"5bbdadf782ebac06a695a8e7"`
	Name     string    `bson:"name" json:"name"`
	Email    string    `bson:"email" json:"email"`
	Password string    `bson:"password" json:"password"`
	Created  time.Time `bson:"created" json:"created"`
}
