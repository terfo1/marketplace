package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"`
	Price       float64            `bson:"price" json:"price"`
	ImageURL    string             `bson:"image_url" json:"image_url"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type Interaction struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	ProductID  primitive.ObjectID `bson:"product_id" json:"product_id"`
	ActionType string             `bson:"action_type" json:"action_type"` // "view" | "like" | "purchase"
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
}
