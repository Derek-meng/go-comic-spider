package topic

import (
	"github.com/Derek-meng/go-comic-spider/dao/topic"
	repo "github.com/Derek-meng/go-comic-spider/repostories/topic"
)

type TopicService struct {
}

func GetLists(page int) []topic.Topic {
	topics := repo.Books(page, 25)
	return topics
}
