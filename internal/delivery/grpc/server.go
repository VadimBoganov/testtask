package grpc

import (
	"context"
	"github.com/VadimBoganov/testtask/internal/services"
	p "github.com/VadimBoganov/testtask/pkg/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"os"
	"time"
)

type Server struct {
	service *services.Service
}

func NewServer(service *services.Service) *Server{
	return &Server{service: service}
}

func (s *Server) FetchFile(req *p.FetchFileRequest, stream p.TestTask_FetchFileServer) error{
	bufferSize := 64 * 1024

	fileName, err := s.service.FetchFile(context.TODO(), req.Url)
	if err != nil{
		return err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	buff := make([]byte, bufferSize)
	for {
		bytesRead, err := file.Read(buff)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		resp := &p.FetchFileResponse{
			FileChunk: buff[:bytesRead],
		}
		err = stream.Send(resp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) GetProducts(ctx context.Context, req *p.ProductsRequest) (*p.ProductsResponse, error){
	products, err := s.service.GetProducts(ctx, req.Paginate.GetLimit(), req.Paginate.Page, req.Sort.GetFieldName(), byte(req.Sort.GetSortType()))
	if err != nil{
		return nil, err
	}

	productResponse := make([]*p.Product, 0, len(products))

	for _, prod := range products{
		productResponse = append(productResponse,
			createProductResponse(prod.ID, prod.Product.Name, prod.Product.Price, prod.PriceChangeCount, prod.LastUpdateTime))
	}

	return &p.ProductsResponse{
		Products: productResponse,
	}, nil
}

func createProductResponse(id primitive.ObjectID, name string, price float64, count int, lastUpdateTime primitive.Timestamp) *p.Product{
	return &p.Product{
		Id: id.Hex(),
		Name: name,
		Price: float32(price),
		PriceChangeCount: int32(count),
		LastUpdateTime: timestamppb.New(time.Unix(int64(lastUpdateTime.T),0)),
	}
}

