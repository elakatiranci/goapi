package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommunityPost struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	UserID    string             `bson:"user_id" json:"user_id"` // Assuming userID is a string
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
}

// The following struct can be used for creating and updating posts
type CommunityPostData struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  string `json:"user_id"`
}
