package main

import (
	"github.com/Derek-meng/go-comic-spider/repostories/topic"
	"github.com/Derek-meng/go-comic-spider/server"
	"github.com/Derek-meng/go-comic-spider/services/spider"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
	_ "github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "comic",
		Usage: "spider comic ",
		Commands: []cli.Command{
			{
				Name:    "server",
				Aliases: []string{"a"},
				Usage:   "start run server",
				Action: func(c *cli.Context) error {

					server.Run()
					return nil
				},
			},
			{
				Name:    "detect",
				Aliases: []string{"a"},
				Usage:   "detect comic",
				Action: func(c *cli.Context) error {
					for _, t := range topic.All() {
						spider.Detector(t)
					}
					return nil
				},
			},
			{
				Name:    "touch",
				Aliases: []string{"a"},
				Usage:   "comic create",
				Action: func(c *cli.Context) error {
					topic.CreateByTitleAndUrl("元尊", "https://www.ohmanhua.com/10136/")
					topic.CreateByTitleAndUrl("凤逆天下", "https://www.ohmanhua.com/10183/")
					topic.CreateByTitleAndUrl("风起苍岚", "https://www.ohmanhua.com/10182/")
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
