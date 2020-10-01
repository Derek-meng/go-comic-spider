package episode_repo

import (
	"context"
	"fmt"
	"github.com/Derek-meng/go-comic-spider/client/db"
	"github.com/Derek-meng/go-comic-spider/dao/episode_dao"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Repo struct {
	mongo   *db.MongoDB
	collect *mongo.Collection
	ctx     context.Context
}

const name = "episodes"

func NewRepo(db *db.MongoDB) Repo {
	return Repo{
		mongo:   db,
		collect: db.DB.Collection(name),
	}
}

func (r Repo) IsExist(e episode_dao.Episode) bool {
	ctx, cancel := defaultCtx()
	err := r.collect.FindOne(ctx, e).Decode(&e)
	defer cancel()
	return err == nil
}

func (r Repo) Create(e episode_dao.Episode) episode_dao.Episode {
	ctx, cancel := defaultCtx()
	defer cancel()
	result, err := r.collect.InsertOne(ctx, e)

	if err != nil {
		log.Fatalf("create episode fail error %s", err)
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); !ok {
		fmt.Println("episode format id fail")
		return episode_dao.Episode{}
	} else {
		e.Id = id
		return e
	}

}

func defaultCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), 300*time.Second)
}

func (r Repo) Get(page, perPage int64, topicId string) []episode_dao.Episode {
	var result []episode_dao.Episode
	opt := options.Find()
	opt.SetLimit(perPage).SetSkip(page).SetSort(bson.D{{"name", -1}})
	ctx, cancelFunc := defaultCtx()
	defer cancelFunc()
	cursor, err := r.collect.Find(ctx, bson.M{"top_id": topicId}, opt)
	if err != nil {
		fmt.Println("get fail", err)
		return []episode_dao.Episode{}
	}
	for cursor.Next(ctx) {
		var ep episode_dao.Episode
		if err := cursor.Decode(&ep); err != nil {
			log.Fatalf("find error :%+v", err)
		}
		result = append(result, ep)
	}
	return result
}

func (r Repo) Count(e episode_dao.Episode) int64 {
	ctx, cancelFunc := defaultCtx()
	defer cancelFunc()
	count, err := r.collect.CountDocuments(ctx, e)
	if err != nil {
		return 0
	}
	return count
}
func (r Repo) Find(e episode_dao.Episode) episode_dao.Episode {
	var ep episode_dao.Episode
	ctx, cancelFunc := defaultCtx()
	defer cancelFunc()

	if err := r.collect.FindOne(ctx, e).Decode(&ep); err != nil {
		return episode_dao.Episode{}
	}
	return ep
}
