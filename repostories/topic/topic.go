package topic

import (
	"fmt"
	"github.com/Derek-meng/go-comic-spider/client/db"
	"github.com/Derek-meng/go-comic-spider/repostories/host"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

const OhComicCode = "OH_COMIC"

type Topic struct {
	Id    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	WebId primitive.ObjectID `json:"web_id,omitempty" bson:"web_id,omitempty"`
	Title string             `json:"title,omitempty" bson:"title,omitempty"`
	Url   string             `json:"url,omitempty" bson:"url,omitempty"`
}

const CollectName = "topic"

func Create() Topic {
	topic := Topic{
		WebId: host.FindByCode().Id,
		Title: "元尊",
		Url:   "https://www.ohmanhua.com/10136/",
	}
	insert, err := getCollection().InsertOne(nil, topic)
	if err != nil {
		log.Fatalf("topic insert fail err: %s", err)
	}
	id, ok := insert.InsertedID.(primitive.ObjectID)
	if !ok {
		fmt.Println(ok)
	}
	topic.Id = id
	return topic

}
func Create2() Topic {
	//凤逆天下  风起苍岚
	//https://www.ohmanhua.com/10183/  https://www.ohmanhua.com/10182/
	topic := Topic{
		WebId: host.FindByCode().Id,
		Title: "凤逆天下",
		Url:   "https://www.ohmanhua.com/10183/",
	}
	insert, err := getCollection().InsertOne(nil, topic)
	if err != nil {
		log.Fatalf("topic insert fail err: %s", err)
	}
	id, ok := insert.InsertedID.(primitive.ObjectID)
	if !ok {
		fmt.Println(ok)
	}
	topic.Id = id
	return topic
}

func getCollection() *mongo.Collection {
	return db.NewDB().Collection(CollectName)
}

func FindByName(name string) Topic {
	var result Topic
	if err := getCollection().FindOne(nil, Topic{Title: name}).Decode(&result); err != nil {
		log.Fatalf("find topic by name:%s error %s\n", name, err)
	}
	return result
}

func All() []Topic {
	var topics []Topic
	cur, err := getCollection().Find(nil, bson.D{})
	if err != nil {
		log.Fatalln("topic find error", err)
	}

	for cur.Next(nil) {
		var t Topic
		err := cur.Decode(&t)
		if err != nil {
			log.Fatal(err)
		}

		topics = append(topics, t)

	}
	return topics
}
