package models

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"` // Note: Store hashed passwords in production!
}