package server

import (
	"fmt"
	"github.com/Derek-meng/go-comic-spider/controller/episode"
	"github.com/Derek-meng/go-comic-spider/controller/topic"
	"github.com/gin-gonic/gin"
	"os"
)

func Run() {

	engine := gin.Default()
	engine.GET("/", topic.Lists)
	engine.GET("episode/", episode.Index)
	engine.GET("episode/info", episode.Info)
	engine.Run(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
