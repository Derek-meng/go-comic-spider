package topic

import (
	"github.com/Derek-meng/go-comic-spider/dao/topic"
	service "github.com/Derek-meng/go-comic-spider/services/topic"
	"net/http"
	"net/url"
	"path"
)

type TopicPresent struct {
	Topic topic.Topic
	URL   string
}

func GetLists(page int, r *http.Request) []TopicPresent {
	var results []TopicPresent
	lists := service.GetLists(page)
	for _, t := range lists {
		parse, err := url.Parse("http://" + r.Host)
		if err != nil {
			return []TopicPresent{}
		}
		q := parse.Query()
		q.Set("id", t.Id.Hex())
		parse.RawQuery = q.Encode()
		parse.Path = path.Join(parse.Path, "episode")
		list := TopicPresent{
			Topic: t,
			URL:   parse.String(),
		}
		results = append(results, list)
	}
	return results

}
