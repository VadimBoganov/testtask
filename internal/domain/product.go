package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	Name string
	Price float64
}

type DBProduct struct {
	ID               primitive.ObjectID
	Product          Product
	PriceChangeCount int
	LastUpdateTime primitive.Timestamp
}
