package grpc

import (
	"context"
	repo "github.com/VadimBoganov/testtask/internal/repository"
	"github.com/VadimBoganov/testtask/internal/services"
	p "github.com/VadimBoganov/testtask/pkg/proto"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMain (m *testing.M)  {
	ctx := context.Background()

	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27018"))

	database := client.Database("test")
	repository := repo.NewRepository(database)
	service := services.NewService(repository)
	server := NewServer(service)

	lis, _ := net.Listen("tcp", ":" + strconv.Itoa(4151))

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	p.RegisterTestTaskServer(grpcServer, server)
	go grpcServer.Serve(lis)

	code := m.Run()
	os.Remove("test_file.csv")
	os.Exit(code)
}

func TestServer_FetchFile_Integration(t *testing.T) {
	if testing.Short(){
		t.Skip()
	}
	conn, _ := grpc.Dial("localhost:" + strconv.Itoa(4151), grpc.WithInsecure())

	c := p.NewTestTaskClient(conn)
	fileStreamResponse, _ := c.FetchFile(context.TODO(), &p.FetchFileRequest{
		Url: "http://localhost:8080/test_file.csv",
	})

	file := make([]byte, 0)
	for {
		chunkResponse, err := fileStreamResponse.Recv()
		if err == io.EOF{
			break
		}
		for _, b := range chunkResponse.FileChunk{
			file = append(file, b)
		}
	}
	assert.NotNil(t, file)
}

func TestServer_GetProducts_Integration(t *testing.T) {
	if testing.Short(){
		t.Skip()
	}
	conn, _ := grpc.Dial("localhost:" + strconv.Itoa(4151), grpc.WithInsecure())

	c := p.NewTestTaskClient(conn)
    c.FetchFile(context.TODO(), &p.FetchFileRequest{
		Url: "http://localhost:8080/test_file.csv",
	})
	time.Sleep(1 * time.Second)
	resp, _ := c.GetProducts(context.TODO(), &p.ProductsRequest{
		Paginate: &p.Paginate{
			Limit: 10,
			Page: 1,
		},
		Sort: &p.Sort{
			FieldName: "product.name",
			SortType: 1,
		},
	})

	assert.NotNil(t, resp.Products)
	assert.Greater(t, len(resp.Products), 0)
}