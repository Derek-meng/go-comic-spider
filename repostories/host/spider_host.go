package host

import (
	"context"
	"github.com/Derek-meng/go-comic-spider/client/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

const SpiderHostCollectName = "spider_host"

type Web struct {
	Id   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	URL  string             `json:"url,omitempty" bson:"url,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Code string             `json:"code,omitempty"  bson:"code,omitempty"`
}

const OhComicCode = "OH_COMIC"

func Insert() Web {
	defer cancel()
	collection := getCollection()
	web := Web{
		URL:  "https://www.ohmanhua.com/",
		Name: "OH 漫畫2",
		Code: OhComicCode,
	}
	one, err := collection.InsertOne(nil, web)
	if err != nil {
		log.Fatalf("insert fail error: %s", err)
	}
	id, ok := one.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Fatalf("spider_host innsert id not assert primitive.ObjectID")
	}
	web.Id = id

	return web

}

var cancel context.CancelFunc

func getCollection() *mongo.Collection {
	var ctx context.Context

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	return db.NewDB(ctx).Collection(SpiderHostCollectName)
}
func FindByCode() Web {

	filter := Web{
		Code: OhComicCode,
	}
	var result Web

	if err := getCollection().FindOne(nil, filter).Decode(&result); err != nil {
		log.Fatalf("decode have error %s", err)
	}
	defer cancel()
	return result

}
