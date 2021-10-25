package configs

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"strings"
	"testing"
)

var config Config

func TestMain(m *testing.M) {
	godotenv.Load("../.env")
	os.Exit(m.Run())
}

func TestConfig_InitConfig_Integration(t *testing.T){
	if testing.Short(){
		t.Skip()
	}

	err := config.InitConfig(".")
	assert.NoErrorf(t, err, "Error occurred while initialization config", err)
}

func TestConfig_MakeMongoDBUri_Integration(t *testing.T) {
	if testing.Short(){
		t.Skip()
	}

	var sb strings.Builder

	config.InitConfig(".")

	sb.WriteString("mongodb://")
	userName, _ := os.LookupEnv("MONGODB_USERNAME")
	sb.WriteString(userName)
	sb.WriteString(":")
	pwd, _ := os.LookupEnv("MONGODB_PASSWORD")
	sb.WriteString(pwd)
	sb.WriteString("@")
	sb.WriteString(config.MongoDB.Endpoints[0].Host)
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(config.MongoDB.Endpoints[0].Port))
	sb.WriteString("/?maxPoolSize=")
	sb.WriteString(strconv.Itoa(config.MongoDB.MaxPoolSize))

	expected := sb.String()

	result := config.MakeMongoDBUri()

	assert.Equal(t, expected, result)
}


