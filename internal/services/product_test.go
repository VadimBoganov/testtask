package services

import (
	"context"
	"github.com/VadimBoganov/testtask/internal/domain"
	"github.com/VadimBoganov/testtask/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"sort"
	"testing"
)

const URL = "http://localhost:8080/"
const FILE_NAME = "test_file.csv"

var productService *Service

func TestMain(m *testing.M){
	ctx := context.Background()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27018"))
	productService = NewService(repository.NewRepository(client.Database("test")))

	code := m.Run()
	os.Remove(FILE_NAME)
	os.Exit(code)
}

func TestProductService_FetchFile_Integration(t *testing.T) {
	if testing.Short(){
		t.Skip()
	}
	fileName, err := productService.FetchFile(context.TODO(), URL + FILE_NAME)
	assert.NoErrorf(t, err, "Error occurred while fetch file", err)
	assert.Greater(t, len(fileName), 0)
}

func TestProductService_GetProducts(t *testing.T) {
	if testing.Short(){
		t.Skip()
	}
	products, err := productService.GetProducts(context.TODO(), 10, 1, "Product.Name", 1)
	assert.NoErrorf(t, err, "Error occurred while get products", err)
	assert.NotNil(t, products)
	assert.Greater(t, len(products), 0)
}

func TestDownloadFile_Integration(t *testing.T) {
	if testing.Short(){
		t.Skip()
	}
	fileName, err := DownloadFile(URL + FILE_NAME)
	assert.NoErrorf(t, err, "Error occurred while download file", err)
	assert.Greater(t, len(fileName), 0)
}

func TestParseCsv_Integration(t *testing.T) {
	if testing.Short(){
		t.Skip()
	}
	DownloadFile(URL + FILE_NAME)
	products, err := ParseCsv(FILE_NAME)
	assert.NoErrorf(t, err, "Error occurred while parse csv file", err)
	assert.NotNil(t, products)
	assert.Greater(t, len(products), 0)
}

func TestMakeDBProducts(t *testing.T) {
	products := []domain.Product{
		{
			Name: "Apple",
			Price: 123,
		},
		{
			Name: "Pineapple",
			Price: 13131,
		},
		{
			Name: "Lemon",
			Price: 3123,
		},
		{
			Name: "Lemon",
			Price: 111,
		},
		{
			Name: "Apple",
			Price: 123,
		},
		{
			Name: "Apple",
			Price: 321311,
		},
	}
	expected := []domain.DBProduct{
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

	result := MakeDBProducts(products)

	sort.Slice(result, func(i, j int) bool {
		return result[i].Product.Name < result[j].Product.Name
	})

	assert.NotNil(t, result)
	assert.Equal(t, expected[0].Product.Name, result[0].Product.Name)
	assert.Equal(t, expected[0].PriceChangeCount, result[0].PriceChangeCount)
}