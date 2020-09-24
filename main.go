package main

import (
	"github.com/Derek-meng/go-comic-spider/repostories/topic"
	"github.com/Derek-meng/go-comic-spider/services/spider"
	_ "github.com/joho/godotenv/autoload"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	for _, t := range topic.All() {
		spider.Detector(t.Url)
	}
}
