package menu

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type SubCategory struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type SubCategorySelector struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetSubCategory(subCategoryURL string, links *[]string) {
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.79 Safari/537.36"))
	c.OnError(func(r *colly.Response, err error) {
		log.Println(err.Error())
	})

	// c.OnResponse(func(r *colly.Response) {
	// 	log.Println(r.StatusCode)
	// })

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("--------------------------------", r.URL, "--------------------------------")
	})

	c.OnHTML("div.subcategories", func(h *colly.HTMLElement) {
		subCategorySelector := SubCategorySelector{
			Name: "a > span",
			URL:  "a",
		}

		h.ForEach("div.category-miniature", func(i int, elm *colly.HTMLElement) {
			// name := (h.DOM.Find(subCategorySelector.Name).Text())
			name := strings.TrimSpace(elm.DOM.Find(subCategorySelector.Name).Text())
			url := elm.DOM.Find(subCategorySelector.URL).AttrOr("href", "-")
			fmt.Println("---> ", name, url)
			*links = append(*links, url)
		})
	})

	c.OnHTML("#main", func(h *colly.HTMLElement) {
		hasSub := h.DOM.ChildrenFiltered("div.subcategories")
		if hasSub.Length() == 0 {
			isEmpty := h.DOM.ChildrenFiltered("#products > #content > div.card > h2")
			if isEmpty.Length() == 0 {
				fmt.Println("============> Scrape products", h.Request.URL.String())
				*links = append(*links, h.Request.URL.String())
			}
		}
	})

	// trigger on empty subcategory
	c.OnHTML("#content", func(h *colly.HTMLElement) {
		txt := h.DOM.Find("div.card > h2").Text()
		fmt.Println("/!\\ ", txt)
	})

	c.OnHTML("#products", func(h *colly.HTMLElement) {
		fmt.Println("+++ A product")
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       1,
		RandomDelay: 1 * time.Second,
	})

	err := c.Visit(subCategoryURL)
	if err != nil {
		log.Println(err)
	}
	// return Category{
	// 	Name: name,
	// 	URL:  url,
	// }
}
