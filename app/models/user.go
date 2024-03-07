package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Login     string             `json:"login" bson:"login"`
	Password  string             `json:"password" bson:"password"`
	Role      string             `json:"role" bson:"role"`
	Issuances []Issuance         `json:"issuances,omitempty" bson:"issuances,omitempty"`
	Token     string             `json:"token" bson:"token"`
	Validated bool               `json:"validated" bson:"validated"`
}

type Issuance struct {
	BookID  string    `json:"book_id" bson:"book_id"`
	DueDate time.Time `json:"due_date" bson:"due_date"`
}
type LoginRequest struct {
	Login    string `json:"login" bson:"login"`
	Password string `json:"password" bson:"password"`
}
