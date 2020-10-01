package episode

import (
	"errors"
	"github.com/Derek-meng/go-comic-spider/presenter/episode"
	service "github.com/Derek-meng/go-comic-spider/services/episode"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

func Index(c *gin.Context) {
	id := c.Query("id")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.Error(errors.New("page error"))
	}
	perPage, err := strconv.Atoi(c.DefaultQuery("page", "50"))
	if err != nil {
		c.Error(errors.New("perPage error"))
	}
	tmpl := template.Must(template.ParseFiles("views/episode/index.html"))
	tmpl.Execute(c.Writer, episode.GetList(id, page, perPage, c.Request))
}

func Info(c *gin.Context) {
	id := c.Query("id")
	tmpl := template.Must(template.ParseFiles("views/episode/info.html"))
	tmpl.Execute(c.Writer, service.NewService().Info(id))
}
