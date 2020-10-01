package topic

import (
	"context"
	"fmt"
	"github.com/Derek-meng/go-comic-spider/client/db"
	"github.com/Derek-meng/go-comic-spider/dao/topic"
	"github.com/Derek-meng/go-comic-spider/repostories/host"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const OhComicCode = "OH_COMIC"

const CollectName = "topic"

type TopicRepo struct {
}

func NewTopicRepo() TopicRepo {
	return TopicRepo{}
}

func CreateByTitleAndUrl(title, url string) topic.Topic {
	t := topic.Topic{
		WebId: host.FindByCode().Id,
		Title: title,
		Url:   url,
	}
	insert, err := getCollection().InsertOne(nil, t)
	if err != nil {
		log.Fatalf("t insert fail err: %s", err)
	}
	id, ok := insert.InsertedID.(primitive.ObjectID)
	if !ok {
		fmt.Println(ok)
	}
	t.Id = id
	return t

}

var cancel context.CancelFunc
var ctx context.Context

func getCollection() *mongo.Collection {

	ctx, cancel = context.WithTimeout(context.TODO(), 10*time.Second)
	return db.NewDB(ctx).Collection(CollectName)
}

func FindByName(name string) topic.Topic {
	var result topic.Topic
	if err := getCollection().FindOne(nil, topic.Topic{Title: name}).Decode(&result); err != nil {
		log.Fatalf("find topic by name:%s error %s\n", name, err)
	}
	return result
}

func Books(page, perPage int) []topic.Topic {
	limit := int64(perPage)
	skip := int64(page - 1)
	option := options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
	}

	cur, err := getCollection().Find(nil, bson.D{}, &option)
	if err != nil {
		log.Fatalln("topic find error", err)

	}
	return decodeLists(cur)
}

func decodeLists(cur *mongo.Cursor) []topic.Topic {
	var result []topic.Topic

	for cur.Next(nil) {
		var t topic.Topic
		err := cur.Decode(&t)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, t)
	}
	return result
}

func All() []topic.Topic {
	var topics []topic.Topic
	cur, err := getCollection().Find(nil, bson.D{})
	defer cancel()
	defer cur.Close(ctx)
	if err != nil {
		log.Fatalln("topic find error", err)
	}

	for cur.Next(nil) {
		var t topic.Topic
		err := cur.Decode(&t)
		if err != nil {
			log.Fatal(err)
		}

		topics = append(topics, t)

	}
	return topics
}
