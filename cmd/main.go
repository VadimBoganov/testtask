package main

import (
	"context"
	g "github.com/VadimBoganov/testtask/internal/delivery/grpc"
	repo "github.com/VadimBoganov/testtask/internal/repository"
	"github.com/VadimBoganov/testtask/internal/services"
	p "github.com/VadimBoganov/testtask/pkg/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"net"
)

const URI = ""

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil{
		//log error
	}

	defer func(){
		if err := client.Disconnect(ctx); err != nil{
			//log error
		}
	}()

	database := client.Database("")

	repository := repo.NewRepository(database)
	service := services.NewService(repository)
	server := g.NewServer(service)

	lis, err := net.Listen("tcp", "")
	if err != nil{
		//log error
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	p.RegisterTestTaskServer(grpcServer, server)
	grpcServer.Serve(lis)
}