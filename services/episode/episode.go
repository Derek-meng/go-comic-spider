package episode

import (
	repo "github.com/Derek-meng/go-comic-spider/repostories/episode"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewService() EpisodeServcie {
	return EpisodeServcie{}
}

type EpisodeServcie struct {
}

func (s EpisodeServcie) GetByTopic(topicId string, page, perPage int64) []repo.Episode {
	objectID, err := primitive.ObjectIDFromHex(topicId)
	if err != nil {
		return []repo.Episode{}
	}
	e := repo.Episode{
		TopicId: objectID,
	}

	return e.Get(page, perPage)
}

func (s EpisodeServcie) Count(topicId string) int64 {
	id, err := primitive.ObjectIDFromHex(topicId)
	if err != nil {
		return 0
	}
	e := repo.Episode{
		TopicId: id,
	}
	return e.Count()
}

func (e EpisodeServcie) Info(id string) repo.Episode {
	dao := repo.NewEpisode()
	object, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return repo.Episode{}
	}
	dao.Id = object
	return dao.Find()
}
