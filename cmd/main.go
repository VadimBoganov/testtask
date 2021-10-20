package main

import (
	"context"
	"github.com/VadimBoganov/testtask/configs"
	g "github.com/VadimBoganov/testtask/internal/delivery/grpc"
	repo "github.com/VadimBoganov/testtask/internal/repository"
	"github.com/VadimBoganov/testtask/internal/services"
	p "github.com/VadimBoganov/testtask/pkg/proto"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil{
		log.Fatalf("Error occured while initialize env variables: %s", err.Error())
	}

	var config configs.Config

	if err := config.InitConfig(); err != nil {
		log.Fatalf("Error occured while initialize config: %s", err.Error())
	}

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MakeMongoDBUri()))
	if err != nil{
		log.Fatalf("Error occured while connect to mongodb: %s", err.Error())
	}

	defer func(){
		if err := client.Disconnect(ctx); err != nil{
			log.Fatalf(err.Error())
		}
	}()

	database := client.Database(viper.GetString("mongodb.databaseName"))

	repository := repo.NewRepository(database)
	service := services.NewService(repository)
	server := g.NewServer(service)

	port := viper.GetInt("grpc-server.port")
	lis, err := net.Listen("tcp", strconv.Itoa(port))
	if err != nil{
		log.Fatalf("Error occured while grpc server listen port: %d  %s", port, err.Error())
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	p.RegisterTestTaskServer(grpcServer, server)
	if err := grpcServer.Serve(lis); err != nil{
		log.Fatalf("Error occured while runnig grpc server: %s", err.Error())
	}
}
