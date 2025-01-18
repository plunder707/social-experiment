// models/post.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Username  string             `bson:"username" json:"username"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt string             `bson:"created_at" json:"created_at"`
}
