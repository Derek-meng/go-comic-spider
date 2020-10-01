package episode

import (
	"errors"
	"github.com/Derek-meng/go-comic-spider/presenter/episode"
	service "github.com/Derek-meng/go-comic-spider/services/episode"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
)

func Index(c *gin.Context) {
	id := c.Query("id")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.Error(errors.New("page error"))
	}
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "50"))
	if err != nil {
		c.Error(errors.New("perPage error"))
	}
	funcMap := template.FuncMap{
		"Iterate": func(count int) []int {
			var Items []int
			for i := 0; i < count; i++ {
				Items = append(Items, i+1)
			}
			return Items
		},
		"Page": func(page int) string {
			parse, _ := url.Parse("http://" + c.Request.Host)
			q := parse.Query()
			q.Set("page", strconv.FormatInt(int64(page), 10))
			q.Set("id", id)
			parse.RawQuery = q.Encode()
			parse.Path = path.Join(parse.Path, "episode")
			return parse.String()
		},
		"IsCurrentPage": func(p int) bool {
			return page == p
		},
	}

	tmpl := template.Must(template.New(filepath.Base("views/episode/index.html")).Funcs(funcMap).ParseFiles("views/episode/index.html"))
	tmpl.Execute(c.Writer, episode.GetList(id, page, perPage, c.Request))
}

func Info(c *gin.Context) {
	id := c.Query("id")
	tmpl := template.Must(template.ParseFiles("views/episode/info.html"))
	tmpl.Execute(c.Writer, service.NewService().Info(id))
}
