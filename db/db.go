package db

import (
	"fmt"
	"os"

	"github.com/couchbase/gocb/v2"
	"github.com/joho/godotenv"
)

func CreateCouchbaseConnection() (*gocb.Cluster, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	connectionString := os.Getenv("COUCHBASE_CONNECTION_STRING")
	dbUsername := os.Getenv("COUCHBASE_USERNAME")
	password := os.Getenv("COUCHBASE_PASSWORD")

	cluster, err := gocb.Connect(connectionString, gocb.ClusterOptions{
		Username: dbUsername,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("error connecting to Couchbase: %w", err)
	}

	return cluster, nil
}
