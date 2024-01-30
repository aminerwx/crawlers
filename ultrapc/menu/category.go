package menu

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

type Category struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type CategorySelector struct {
	NameSelector string `json:"name"`
	URLSelector  string `json:"url"`
}

func GetCategory(currentMenuUrl string, links *[]string) {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.79 Safari/537.36"))
	c.OnError(func(r *colly.Response, err error) {
		log.Println(err.Error())
	})

	// c.OnResponse(func(r *colly.Response) {
	// 	log.Println(r.StatusCode)
	// })

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("#######################", r.URL, "#######################")
	})

	c.OnHTML("#main", func(h *colly.HTMLElement) {
		CategorySelector := CategorySelector{
			NameSelector: "a > span",
			URLSelector:  "a",
		}

		h.ForEach("div.subcategories > div.category-miniature", func(i int, elm *colly.HTMLElement) {
			// name := strings.TrimSpace(h.DOM.Find(subCategorySelector.Name).Text())
			url := elm.DOM.Find(CategorySelector.URLSelector).AttrOr("href", "-")
			// fmt.Println("@@@@@@@@@@@", url)
			GetSubCategory(url, links)
		})
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       1,
		RandomDelay: 1 * time.Second,
	})

	err := c.Visit(currentMenuUrl)
	if err != nil {
		log.Println(err)
	}
	// return Category{
	// 	Name: name,
	// 	URL:  url,
	// }
}
