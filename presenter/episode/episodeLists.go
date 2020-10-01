package episode

import (
	service "github.com/Derek-meng/go-comic-spider/services/episode"
	"net/http"
	"net/url"
	"path"
)

type episode struct {
	Id   string
	Name string
	Url  string
}
type Lists struct {
	Episodes    []episode
	PageCount   int
	CurrentPage int
}

func GetList(id string, page, perPage int, r *http.Request) Lists {
	var result []episode
	s := service.NewService()
	lists := s.GetByTopic(id, int64(page), int64(perPage))
	c := s.Count(id)
	for _, list := range lists {
		parse, err := url.Parse("http://" + r.Host)
		if err != nil {
			return Lists{}
		}
		q := parse.Query()
		q.Set("id", list.Id.Hex())
		parse.RawQuery = q.Encode()
		parse.Path = path.Join(parse.Path, "episode", "info")
		var e = episode{
			Id:   list.Id.Hex(),
			Name: list.Name,
			Url:  parse.String(),
		}
		result = append(result, e)
	}
	count := int(c) / perPage
	if (int(c) % perPage) > 0 {
		count++
	}
	return Lists{
		Episodes:    result,
		PageCount:   count,
		CurrentPage: page,
	}
}
