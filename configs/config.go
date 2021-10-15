package configs

import (
	"os"
	"strings"
)

type Config struct {
	mongoDB MongoDB
	csv Csv
	grpcServer GrpcServer
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

	endpoints := c.mongoDB.endpoints

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
	sb.WriteString(string(c.mongoDB.maxPoolSize))

	if len(endpoints) > 1{
		sb.WriteString("&replicaSet=RS")
	}

	return sb.String()
}

func addEndpoint(endpoint Endpoint, sb *strings.Builder) {
	sb.WriteString(endpoint.host)
	sb.WriteString(":")
	sb.WriteString(string(endpoint.port))
}

type MongoDB struct {
	endpoints []Endpoint
	maxPoolSize int
	databaseName string
}

type Endpoint struct {
	host string
	port int
}

type Csv struct {
	fileName string
}

type GrpcServer struct {
	port int
}