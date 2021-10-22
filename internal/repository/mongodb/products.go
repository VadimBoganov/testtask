package mongodb

import (
	"context"
	"github.com/VadimBoganov/testtask/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Sort struct{
	fieldName string
	sortType int
}

type Paginate struct {
	limit int64
	page int64
}

func newSort(fieldName string, sortType int) *Sort{
	return &Sort{
		fieldName: fieldName,
		sortType: sortType,
	}
}

func newPaginate(limit, page int) *Paginate {
	return &Paginate{
		limit: int64(limit),
		page: int64(page),
	}
}

func (p *Paginate) setPaginateOpts() *options.FindOptions{
	limit := p.limit
	skip := p.page * p.limit - p.limit
	return &options.FindOptions{Limit: &limit, Skip: &skip}
}

func (s *Sort) setSortOpts(opts *options.FindOptions){
	opts.SetSort(bson.D{{s.fieldName, s.sortType}})
}

type ProductsRepo struct {
	collection *mongo.Collection
}

func NewProductsRepo(db *mongo.Database) *ProductsRepo{
	return &ProductsRepo{
		collection: db.Collection("productsCollection"),
	}
}

func (r *ProductsRepo) Insert(ctx context.Context, products []domain.DBProduct) error {
	var i []interface{}
	for _, prod := range products {
		i = append(i, prod)
	}

	_, err := r.collection.InsertMany(ctx, i)
	return err
}

func (r *ProductsRepo) Get(ctx context.Context, limit, page int, fieldName string, sortType int) ([]domain.DBProduct, error) {
	result := make([]domain.DBProduct, 0)

	opts := newPaginate(limit, page).setPaginateOpts()
	newSort(fieldName, sortType).setSortOpts(opts)

	curr, err := r.collection.Find(ctx, bson.D{}, opts)

	if err != nil{
		return result, err
	}

	for curr.Next(ctx) {
		var elem domain.DBProduct
		if err := curr.Decode(&elem); err != nil{
			return nil, err
		}

		result = append(result, elem)
	}

	return result, nil
}

