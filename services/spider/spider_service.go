package spider

import (
	"fmt"
	"github.com/Derek-meng/go-comic-spider/dao/episode_dao"
	"github.com/Derek-meng/go-comic-spider/dao/topic"
	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

type Repo interface {
	Create(e episode_dao.Episode) episode_dao.Episode
	IsExist(e episode_dao.Episode) bool
}
type Service struct {
	repo Repo
}

func NewService(r Repo) Service {
	return Service{
		repo: r,
	}
}

func (s Service) getImages(u episode_dao.Episode, t int) ([]string, int) {
	t++
	driver, page, err := s.newPage()
	if err != nil {
		return []string{}, t
	}
	defer driver.Stop()
	defer page.CloseWindow()
	if err := page.Navigate(u.Url); err != nil {
		return []string{}, t
	}
	time.Sleep(100 * time.Millisecond)
	pageClass := page.FindByID("mangalist")
	_ = pageClass.ScrollFinger(0, 1700)
	i := 1
	images := make([]string, 0, 10)
	for {
		selection := pageClass.FirstByXPath("//*[@id=\"mangalist\"]/div[" + strconv.Itoa(i) + "]/img")
		elements, err := selection.Elements()
		if err != nil {
			break
		}
		element := elements[0]
		x, y, err := element.GetLocation()
		if err != nil {
			break
		}
		_ = selection.ScrollFinger(x, y)
		attribute, err := element.GetAttribute("src")
		if err != nil {
			break
		}
		images = append(images, attribute)
		i++
	}
	return images, t
}

func (s Service) newPage() (*agouti.WebDriver, *agouti.Page, error) {
	diver := s.newDiver()
	page, err := diver.NewPage()
	if err != nil {
		return diver, page, err
	}
	return diver, page, err
}

func (Service) newDiver() *agouti.WebDriver {
	driver := agouti.ChromeDriver(agouti.ChromeOptions("args", optionList))
	driver.Timeout = 10 * time.Second
	if err := driver.Start(); err != nil {
		log.Fatal("Failed to start driver:", err)
	}
	return driver
}

func (s Service) Detector(t topic.Topic, isBreak bool) {
	host, err := url.Parse(t.Url)
	if err != nil {
		log.Fatalln(err)
	}
	response := s.getResponse(host)
	defer response.Body.Close()

	eps := s.parseHtml(response, host)
	channel := make(chan episode_dao.Episode, 10)
	count := os.Getenv("CHROME_COUNT")
	c, err := strconv.Atoi(count)
	if err != nil {
		log.Fatalln("env CHROME_COUNT is empty")
	}
	var wg sync.WaitGroup
	wg.Add(c)
	for i := 0; i < c; i++ {
		go func() {
		Test:
			for {
				select {
				case e, ok := <-channel:
					var (
						i      = 0
						images []string
					)
					for len(images) == 0 && i < 3 {
						images, i = s.getImages(e, i)
					}
					if len(images) > 0 {
						e.Images = images
						e.TopicId = t.Id
						s.repo.Create(e)
					}
					if !ok {
						break Test
					}

				}
			}
			wg.Done()
		}()
	}

	for _, ep := range eps {
		if s.repo.IsExist(ep) {
			if isBreak {
				break
			} else {
				continue
			}

		} else {
			channel <- ep
		}
	}
	close(channel)
	wg.Wait()
}

func (Service) getResponse(host *url.URL) *http.Response {
	res, err := http.Get(host.String())
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	return res
}

func (Service) parseHtml(res *http.Response, host *url.URL) []episode_dao.Episode {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	class := "body > div.fed-main-info.fed-min-width > div > div.fed-tabs-info.fed-rage-foot.fed-part-rows.fed-part-layout.fed-back-whits.fed-play-data > div > div.fed-tabs-item.fed-drop-info.fed-visible > div.fed-drop-boxs.fed-drop-btms.fed-matp-v > div.fed-play-item.fed-drop-item.fed-visible > div > ul > li"
	eps := make([]episode_dao.Episode, 0, 100)
	doc.Selection.Find(class).Each(func(i int, selection *goquery.Selection) {
		uri, exists := selection.Children().Attr("href")
		if exists {
			title := selection.Children().Text()
			e := episode_dao.Episode{
				Name: title,
				Url:  fmt.Sprintf("%s://%s/%s", host.Scheme, host.Host, uri),
			}
			eps = append(eps, e)
		}
	})
	return eps
}
