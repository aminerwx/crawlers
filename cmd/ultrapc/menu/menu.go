package menu

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocolly/colly"
)

type Menu struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type MenuSelector struct {
	NameSelector string `json:"name"`
	URLSelector  string `json:"url"`
}

type Link struct {
	Name        string
	URL         string
	CategoryURL string
	MenuURL     string
}

func GetMenu() []string {
	var links []string
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.79 Safari/537.36"))

	c.OnError(func(r *colly.Response, err error) {
		log.Println(err.Error())
	})

	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	c.Async = true

	// c.OnResponse(func(r *colly.Response) {
	// 	log.Println(r.StatusCode)
	// })

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("# ", r.URL)
	})

	c.OnHTML("#top-menu", func(h *colly.HTMLElement) {
		menuSelector := []MenuSelector{
			// Components
			{
				NameSelector: "#category-20 > a",
				URLSelector:  "#category-20 > a",
			},
			// Peripherals
			{
				NameSelector: "#category-58 > a",
				URLSelector:  "#category-58 > a",
			},
		}

		// name := strings.TrimSpace(h.DOM.Find(subCategorySelector.Name).Text())
		url := h.DOM.Find(menuSelector[0].URLSelector).AttrOr("href", "-")
		GetCategory(url, &links)

	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       1,
		RandomDelay: 1 * time.Second,
	})

	err := c.Visit("https://www.ultrapc.ma/")
	if err != nil {
		log.Println(err)
	}
	c.Wait()
	// return Category{
	// 	Name: name,
	// 	URL:  url,
	// }

	return links
}
