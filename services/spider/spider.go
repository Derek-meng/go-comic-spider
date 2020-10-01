package spider

import (
	"fmt"
	"github.com/Derek-meng/go-comic-spider/dao/topic"
	"github.com/Derek-meng/go-comic-spider/repostories/episode"
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

var driver *agouti.WebDriver

var optionList = []string{
	"start-maximized",
	"--headless",
	"enable-automation",
	"--window-size=1000,900",
	"--incognito", //隐身模式
	"--blink-settings=imagesEnabled=true",
	"--no-default-browser-check",
	"--ignore-ssl-errors=true",
	"--ssl-protocol=any",
	"--no-sandbox",
	"--disable-breakpad",
	"--disable-logging",
	"--no-zygote",
	"--allow-running-insecure-content",
	"--disable-extensions",
	"--disable-infobars",
	"--disable-dev-shm-usage",
	"--disable-cache",
	"--disable-application-cache",
	"--disable-offline-load-stale-cache",
	"--disk-cache-size=0",
	"--disable-gpu",
	"--dns-prefetch-disable",
	"--no-proxy-server",
	"--silent",
	"--disable-browser-side-navigation",
}

func init() {
	driver = agouti.ChromeDriver(agouti.ChromeOptions("args", optionList))
	if err := driver.Start(); err != nil {
		log.Fatal("Failed to start driver:", err)
	}
}

func getImages(u episode.Episode, t int) ([]string, int) {
	t++
	page := newPage()
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

func newPage() *agouti.Page {
	page, err := driver.NewPage()
	if err != nil {
		log.Fatal("Failed to open newPage:", err)
	}
	return page
}

func Detector(t topic.Topic) {
	host, err2 := url.Parse(t.Url)
	if err2 != nil {
		log.Fatalln(err2)
	}
	res, err := http.Get(host.String())
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	eps := parseHtml(err, res, host)
	channel := make(chan episode.Episode, 10)
	count := os.Getenv("CHROME_COUNT")
	var c int
	c, err = strconv.Atoi(count)
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
					images, i := getImages(e, 1)
					for len(images) == 0 && i < 3 {
						images, i = getImages(e, i)
					}
					if len(images) > 0 {
						e.Images = images
						e.TopicId = t.Id
						e.Create()
					}
					if !ok {
						break Test
					}
				}
			}
			driver.Stop()
			wg.Done()
		}()
	}

	for _, ep := range eps {
		if ep.IsExistsByNameAndURL() {
			break
		} else {
			channel <- ep
		}
	}
	close(channel)
	wg.Wait()
}

func parseHtml(err error, res *http.Response, host *url.URL) []episode.Episode {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	class := "body > div.fed-main-info.fed-min-width > div > div.fed-tabs-info.fed-rage-foot.fed-part-rows.fed-part-layout.fed-back-whits.fed-play-data > div > div.fed-tabs-item.fed-drop-info.fed-visible > div.fed-drop-boxs.fed-drop-btms.fed-matp-v > div.fed-play-item.fed-drop-item.fed-visible > div > ul > li"
	eps := make([]episode.Episode, 0, 100)
	doc.Selection.Find(class).Each(func(i int, selection *goquery.Selection) {
		uri, exists := selection.Children().Attr("href")
		if exists {
			title := selection.Children().Text()
			e := episode.Episode{
				Name: title,
				Url:  fmt.Sprintf("%s://%s/%s", host.Scheme, host.Host, uri),
			}
			eps = append(eps, e)
		}
	})
	return eps
}
