package services

import (
	"context"
	"encoding/csv"
	"github.com/VadimBoganov/testtask/internal/domain"
	"github.com/VadimBoganov/testtask/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	u "net/url"
	"os"
	"strconv"
	"time"
)

type ProductService struct {
	repo repository.Product
}

func NewProductService(repo repository.Product) *ProductService{
	return &ProductService{
		repo: repo,
	}
}

// FetchFile TODO:make depends for file downloader and csv parser(refactoring)
func (s *ProductService) FetchFile(ctx context.Context, url string) error{
	parsedUrl, err := u.Parse(url)
	fileName := parsedUrl.Path[1:]

	DownloadFile(url, fileName)

	products, err := ParseCsv(fileName)
	if err != nil{
		return err
	}

	dbProducts := makeDBProducts(products)
	err = s.repo.Insert(ctx, dbProducts)

	return err
}

func (s *ProductService) GetProducts(ctx context.Context, limit, page int32, fieldName string, sortType byte) ([]domain.DBProduct, error){
	return s.repo.Get(ctx, int(limit), int(page), fieldName, sortType)
}

func DownloadFile(url, fileName string) error{
	resp, err := http.Get(url)
	if err != nil{
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(fileName)
	if err != nil{
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return  err
}

func ParseCsv(path string) ([]domain.Product, error){
	file, err := os.Open(path)
	if err != nil{
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	products := make([]domain.Product, 0)
	for {
		row, err := csvReader.Read()
		if err != nil{
			if err == io.EOF{
				err = nil
			}
			return products, err
		}

		price, err := strconv.ParseFloat(row[1], 64)
		if err != nil{
			return nil, err
		}

		product := domain.Product{
			Name: row[0],
			Price: price,
		}

		products = append(products, product)
	}
}

func makeDBProducts(products []domain.Product) []domain.DBProduct{
	dbProducts := make(map[string]domain.DBProduct)

	count := 0
	for count < len(products) {
		currProduct := products[count]
		if product, ok := dbProducts[currProduct.Name]; ok{
			if product.Product.Price != currProduct.Price{
				product.PriceChangeCount++
				product.Product.Price = currProduct.Price
				dbProducts[currProduct.Name] = product
			}
		} else {
			product := products[count]
			dbProduct := domain.DBProduct{
				Product:          product,
				PriceChangeCount: 0,
				LastUpdateTime:   primitive.Timestamp{T: uint32(time.Now().Unix()), I: 1},
			}
			dbProducts[products[count].Name] = dbProduct
		}
		count++
	}

	result := make([]domain.DBProduct, 0, len(dbProducts))
	for _, val := range dbProducts{
		result = append(result, val)
	}

	return result
}
