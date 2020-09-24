package episode

import (
	"github.com/Derek-meng/go-comic-spider/client/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Episode struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TopicId primitive.ObjectID `json:"top_id,omitempty" bson:"top_id,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Images  []string           `json:"images,omitempty" bson:"image,omitempty"`
	Url     string             `json:"url,omitempty" bson:"url,omitempty"`
}

const collectName = "episodes"

func getCollection() *mongo.Collection {
	return db.NewDB().Collection(collectName)
}

func (e Episode) IsExistsByNameAndURL() bool {
	err := getCollection().FindOne(nil, e).Decode(&e)
	return err == nil
}

func (e *Episode) Create() {
	result, err := getCollection().InsertOne(nil, e)
	if err != nil {
		log.Fatalf("create episode fail error %s", err)
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); !ok {
		log.Fatalln("primitive.ObjectID fail")
	} else {
		e.Id = id
	}
}
