package main

import (
	"context"
	"fmt"
	"github.com/Derek-meng/go-comic-spider/client/db"
	"github.com/Derek-meng/go-comic-spider/repostories/episode_repo"
	"github.com/Derek-meng/go-comic-spider/repostories/topic"
	"github.com/Derek-meng/go-comic-spider/server"
	"github.com/Derek-meng/go-comic-spider/services/spider"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
)

func main() {
	app := &cli.App{
		Name:  "comic",
		Usage: "spider comic ",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "username,u",
				Usage: "user account",
			},
			cli.StringFlag{
				Name:  "password,p",
				Usage: "user password",
			},
		},
		Commands: []cli.Command{
			{
				Name:  "server",
				Usage: "start run server",
				Action: func(c *cli.Context) error {
					server.Run()
					return nil
				},
			},
			{
				Name:  "detect",
				Usage: "detect comic",
				Action: func(c *cli.Context) error {
					ctx, cancelFunc := context.WithCancel(context.Background())
					defer cancelFunc()
					s := spider.NewService(episode_repo.NewRepo(db.Instance(ctx)))
					for _, t := range topic.All() {
						s.Detector(t, true)
					}
					return nil
				},
			},
			{
				Name:  "fill",
				Usage: "detect comic",
				Action: func(c *cli.Context) error {
					ctx, cancelFunc := context.WithCancel(context.Background())
					defer cancelFunc()
					s := spider.NewService(episode_repo.NewRepo(db.Instance(ctx)))
					for _, t := range topic.All() {
						s.Detector(t, false)
					}
					return nil
				},
			},
			{
				Name:  "touch",
				Usage: "comic create",
				Action: func(c *cli.Context) error {
					topic.CreateByTitleAndUrl("元尊", "https://www.ohmanhua.com/10136/")
					topic.CreateByTitleAndUrl("凤逆天下", "https://www.ohmanhua.com/10183/")
					topic.CreateByTitleAndUrl("风起苍岚", "https://www.ohmanhua.com/10182/")
					return nil
				},
			},
			{
				Name:  "kill",
				Usage: "kill chrome",
				Action: func(c *cli.Context) error {
					out, err := exec.Command("sudo -S pkill -SIGINT chrome").Output()
					fmt.Printf("combined out:\n%s\n", string(out))
					return err
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
