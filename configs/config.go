package configs

import (
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	MongoDB MongoDB `mapstructure:"mongodb"`
	GrpcServer GrpcServer `mapstructure:"grpc_server"`
}

type MongoDB struct {
	Endpoints []Endpoint
	MaxPoolSize int `mapstructure:"max_pool_size"`
	DatabaseName string `mapstructure:"database_name"`
}

type Endpoint struct {
	Host string
	Port int
}

type GrpcServer struct {
	Port int
}

func (c *Config) InitConfig() error{
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil{
		return err
	}

	return viper.Unmarshal(&c)
}

func (c *Config) MakeMongoDBUri() string{
	var sb strings.Builder
	sb.WriteString("mongodb://")

	if val, exists := os.LookupEnv("MONGODB_USERNAME"); exists{
		sb.WriteString(val)
		sb.WriteString(":")
	}

	if val, exists := os.LookupEnv("MONGODB_PASSWORD"); exists{
		sb.WriteString(val)
		sb.WriteString("@")
	}

	endpoints := c.MongoDB.Endpoints

	if len(endpoints) > 1 {
		for _, endpoint := range endpoints {
			addEndpoint(endpoint, &sb)
			sb.WriteString(",")
		}
		currString := sb.String()
		currString = currString[:len(currString)-1]
		sb.Reset()
		sb.WriteString(currString)
	} else {
		addEndpoint(endpoints[0], &sb)
	}

	sb.WriteString("/?maxPoolSize=")
	sb.WriteString(strconv.Itoa(c.MongoDB.MaxPoolSize))

	if len(endpoints) > 1{
		sb.WriteString("&replicaSet=RS")
	}

	return sb.String()
}

func addEndpoint(endpoint Endpoint, sb *strings.Builder) {
	sb.WriteString(endpoint.Host)
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(endpoint.Port))
}

