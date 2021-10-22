package mongodb_test

import (
	"context"
	"github.com/VadimBoganov/testtask/internal/domain"
	"github.com/VadimBoganov/testtask/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

var dbProducts []domain.DBProduct
var repo *repository.Repository

func TestMain(m *testing.M){
	dbProducts = []domain.DBProduct{
		{
			Product: domain.Product{
				Name: "Apple",
				Price: 321311,
			},
			PriceChangeCount: 1,
		},
		{
			Product: domain.Product{
				Name: "Pineapple",
				Price: 13131,
			},
			PriceChangeCount: 0,
		},
		{
			Product: domain.Product{
				Name: "Lemon",
				Price: 111,
			},
			PriceChangeCount: 1,
		},
	}

	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27018"))
	repo = repository.NewRepository(client.Database("test"))

	os.Exit(m.Run())
}

func TestProductsRepo_Insert_Integration(t *testing.T) {
	if testing.Short(){
		t.Skip()
	}
	err := repo.Insert(context.TODO(), dbProducts)
	assert.NoErrorf(t, err, "Error occurred while insert products", err)
}

func TestProductsRepo_Get_Integration(t *testing.T) {
	if testing.Short(){
		t.Skip()
	}

	ctx := context.TODO()
	repo.Insert(ctx, dbProducts)

	result, _ := repo.Get(ctx, 10, 1, "product.name", 1)

	assert.NotNil(t, result)
	assert.Greater(t, len(result), 0)
}

func TestProductsRepo_Get_Limit_Integration(t *testing.T) {
	if testing.Short(){
		t.Skip()
	}

	ctx := context.TODO()
	repo.Insert(ctx, dbProducts)

	result, _ := repo.Get(ctx, 1, 1, "product.name", 1)

	assert.NotNil(t, result)
	assert.Equal(t, len(result), 1)
}

func TestProductsRepo_Get_Page_Integration(t *testing.T) {
	if testing.Short(){
		t.Skip()
	}

	ctx := context.TODO()
	repo.Insert(ctx, dbProducts)

	result, _ := repo.Get(ctx, 1, 3, "product.name", -1)

	assert.NotNil(t, result)
	assert.Equal(t, len(result), 1)
	assert.Equal(t, "Pineapple", result[0].Product.Name)
}

func TestProductsRepo_Get_Sort_Integration(t *testing.T) {
	if testing.Short(){
		t.Skip()
	}

	ctx := context.TODO()
	repo.Insert(ctx, dbProducts)

	result, _ := repo.Get(ctx, 10, 1, "product.name", -1)

	assert.NotNil(t, result)
	assert.Equal(t, "Pineapple", result[0].Product.Name)
}