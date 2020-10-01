package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"sync"
)

type MongoDB struct {
	DB      *mongo.Database
	Session *mongo.Client
}

var instance *MongoDB
var once sync.Once

func Instance(ctx context.Context) *MongoDB {
	once.Do(func() {
		option := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT")))
		session, err := mongo.NewClient(option)
		if err != nil {
			log.Fatalf("mongo connecnt fail %e\n", err)
		}
		if err := session.Connect(ctx); err != nil {
			log.Fatalf("session connecnt fail %e\n", err)
		}

		db := MongoDB{
			Session: session,
			DB:      session.Database("comic"),
		}
		instance = &db
	})
	return instance
}

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
