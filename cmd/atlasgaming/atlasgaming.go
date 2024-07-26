package atlasgaming

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Product struct {
	Name         string
	Link         string
	Menu         string
	MenuLink     string
	Category     string
	CategoryLink string
	Platform     string
	Price        int
	Stock        int
	CurrentPrice int
	Discount     int
}

func Runner() {
}

func GetCategories() []string {
	var categories []string
	categoryCollector := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"))
	url := "https://atlasgaming.ma"
	categoryCollector.OnHTML("li", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "menufirstchild menu-item menu-item-type-taxonomy menu-item-object-product_cat") {
			categories = append(categories, e.ChildAttr("a", "href"))
		}
	})
	categoryCollector.Visit(url)
	categoryCollector.Wait()
	return categories
}

func Crawl(categoryURL string) {
	c := colly.NewCollector(colly.Async(true), colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"))
	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.Request.URL.String() + "RES")
		if r.Request.URL.String() == "https://atlasgaming.ma" {
			fmt.Println("REDIRECTED: ", r.Request.URL.String())
			return
		}
	})

	c.OnHTML("div", func(h *colly.HTMLElement) {
		if strings.Contains(h.Attr("class"), "ct-div-block product-card__info pad--xs") {
			fmt.Println(h.DOM.Find("a > h1").Text())
		}
	})

	for i := 1; i < 100; i++ {
		if i == 3 {
			break
		}
		url := categoryURL + "page/" + strconv.Itoa(i)
		fmt.Println(url)
		c.Visit(url)
		c.Wait()
	}
}
