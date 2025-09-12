package db

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func NewDriver() neo4j.DriverWithContext {
	// Service name/port from docker-compose below
	uri := "bolt://followers-db:7687"
	auth := neo4j.BasicAuth("neo4j", "followersPass", "")
	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	return driver
}

func CloseDriver(ctx context.Context, driver neo4j.DriverWithContext) {
	if err := driver.Close(ctx); err != nil {
		log.Printf("Error closing Neo4j driver: %v", err)
	}
}
