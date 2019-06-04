package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name"`
	Completed bool               `json:"completed"`
	Due       time.Time          `json:"due,omitempty"`
}

type Todos []Todo
