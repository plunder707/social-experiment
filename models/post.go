// models/post.go
package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
    ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    UserID    primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
    Username  string             `json:"username,omitempty" bson:"username,omitempty"`
    Content   string             `json:"content,omitempty" bson:"content,omitempty"`
    CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
