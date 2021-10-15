package services

import (
	"context"
	"github.com/VadimBoganov/testtask/internal/domain"
	"github.com/VadimBoganov/testtask/internal/repository"
)

type Product interface {
	FetchFile(ctx context.Context, url string) (string, error)
	GetProducts(ctx context.Context, limit, page int32, fieldName string, sortType byte) ([]domain.DBProduct, error)
}

type Service struct {
	Product
}

func NewService(repo *repository.Repository) *Service{
	return &Service{
		Product: NewProductService(repo.Product),
	}
}