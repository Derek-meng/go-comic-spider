package episode

import (
	"context"
	"fmt"
	"github.com/Derek-meng/go-comic-spider/client/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Episode struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TopicId primitive.ObjectID `json:"top_id,omitempty" bson:"top_id,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Images  []string           `json:"images,omitempty" bson:"image,omitempty"`
	Url     string             `json:"url,omitempty" bson:"url,omitempty"`
}

const collectName = "episodes"

var cancel context.CancelFunc
var ctx context.Context

func getCollection() *mongo.Collection {

	ctx, cancel = context.WithTimeout(context.TODO(), 10*time.Second)

	return db.NewDB(ctx).Collection(collectName)
}

func NewEpisode() Episode {
	return Episode{}
}

func (e Episode) IsExistsByNameAndURL() bool {
	err := getCollection().FindOne(ctx, e).Decode(&e)
	defer cancel()
	return err == nil
}

func (e *Episode) Create() {
	result, err := getCollection().InsertOne(ctx, e)
	defer cancel()

	if err != nil {
		log.Fatalf("create episode fail error %s", err)
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); !ok {
		log.Fatalln("primitive.ObjectID fail")
	} else {
		e.Id = id
	}
}

func (e Episode) Get(page, perPage int64) []Episode {
	var result []Episode

	opt := options.Find()
	opt.SetLimit(perPage).SetSkip(page).SetSort(bson.D{{"name", -1}})
	cursor, err := getCollection().Find(ctx, e, opt)
	defer cancel()
	if err != nil {
		fmt.Println(err)
		return []Episode{}
	}
	for cursor.Next(ctx) {
		var ep Episode
		if err := cursor.Decode(&ep); err != nil {
			log.Fatalf("find error :%+v", err)
		}
		result = append(result, ep)
	}
	return result
}
func (e Episode) Count() int64 {
	count, err := getCollection().CountDocuments(ctx, e)
	if err != nil {
		return 0
	}
	return count
}
func (e Episode) Find() Episode {
	var ep Episode
	if err := getCollection().FindOne(ctx, e).Decode(&ep); err != nil {
		return Episode{}
	}
	return ep
}
