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
)

func main() {
	if err := initEnv(); err != nil{
		log.Fatalf("Error occured while initialize env variables: %s", err.Error())
	}

	var config configs.Config

	if err := initConfig(&config); err != nil {
		log.Fatalf("Error occured while initialize config: %s", err.Error())
	}

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString(config.MakeMongoDBUri())))
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
	lis, err := net.Listen("tcp", string(port))
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

func initConfig(config *configs.Config) error{
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	if err := viper.ReadInConfig(); err != nil{
		return err
	}

	return viper.Unmarshal(&config)
}

func initEnv() error{
	return godotenv.Load()
}
