package repository

import (
	"context"
	"github.com/VadimBoganov/testtask/internal/domain"
	"github.com/VadimBoganov/testtask/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Product interface {
	Get(ctx context.Context, limit, page int, fieldName string, sortType int) ([]domain.DBProduct, error)
	Insert(ctx context.Context, products []domain.DBProduct) error
}

type Repository struct {
	Product
}

func NewRepository(db *mongo.Database) *Repository{
	return &Repository{
		Product: mongodb.NewProductsRepo(db),
	}
}