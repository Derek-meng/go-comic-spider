package topic

import (
	presenter "github.com/Derek-meng/go-comic-spider/presenter/topic"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

func Lists(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	p, err := strconv.Atoi(page)
	if err != nil {
		p = 1
	}
	tmpl := template.Must(template.ParseFiles("views/index.html"))
	lists := presenter.GetLists(p, c.Request)
	tmpl.Execute(c.Writer, lists)
}
