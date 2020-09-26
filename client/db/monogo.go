package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func NewDB(ctx context.Context) *mongo.Database {
	option := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT")))
	client, err := mongo.NewClient(option)
	if err != nil {
		log.Fatalf("mongo connecnt fail %e\n", err)
	}
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("client connecnt fail %e\n", err)
	}

	return client.Database("comic")
}
