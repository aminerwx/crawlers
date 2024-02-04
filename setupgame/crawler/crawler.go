package crawler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aminerwx/crawlers/helper"
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

func Crawler(url string, products *[]Product) {
	c := colly.NewCollector(colly.Async(true), colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"))
	pages := 0
	totalPages(c, &pages)
	fmt.Println("crawler --> ", url)
	c.Visit(url)
	c.Wait()
	for i := 1; i < pages+1; i++ {
		var sb strings.Builder
		sb.WriteString(url)
		sb.WriteString("?jsf=jet-engine:shop_category&pagenum=")
		sb.WriteString(strconv.Itoa(i))
		getProducts(c, sb.String(), products)
		if i == 1 {
			break
		}
	}
}

func totalPages(c *colly.Collector, pages *int) {
	c.OnHTML(".jet-engine-query-count", func(element *colly.HTMLElement) {
		items, _ := strconv.Atoi(element.Text)
		if items > 0 {
			*pages = (items / 16) + 1
		}
	})
}

func getProducts(c *colly.Collector, url string, products *[]Product) {
	productListSelector := ".jet-listing-grid"
	d := c.Clone()
	d.OnHTML(productListSelector, func(element *colly.HTMLElement) {
		breadcrumb := helper.ParseLink(url)
		element.ForEach("div.jet-listing-grid__items > div.jet-listing-grid__item", func(_ int, product *colly.HTMLElement) {
			productPriceSelector := "div:nth-child(1) > a:nth-child(1) > div:nth-child(4) > div:nth-child(1) > div:nth-child(1) > span:nth-child(1)"
			productCurrentPriceSelector := "div:nth-child(1) > a:nth-child(1) > div:nth-child(4) > div:nth-child(2) > div:nth-child(1) > span:nth-child(1)"
			productName := product.DOM.Find("div:nth-child(1) > a.elementor-element > div > div > div.jet-listing-dynamic-field__content").Text()
			productLink := product.DOM.Find("div:nth-child(1) > a.elementor-element").AttrOr("href", "-")
			productPriceStr := strings.ReplaceAll(product.DOM.Find(productPriceSelector).Text(), ",00 MAD", "")
			productCurrentPriceStr := strings.ReplaceAll(product.DOM.Find(productCurrentPriceSelector).Text(), ",00 MAD", "")
			if productCurrentPriceStr == "" {
				productCurrentPriceStr = "0"
			}
			productPrice, _ := strconv.Atoi(productPriceStr)
			productCurrentPrice, _ := strconv.Atoi(productCurrentPriceStr)
			discount := 0
			if productCurrentPrice > 0 {
				discount = productPrice - productCurrentPrice
			}
			*products = append(*products, Product{
				Name:         productName,
				Link:         productLink,
				Menu:         breadcrumb[0],
				MenuLink:     "/" + breadcrumb[0],
				Category:     breadcrumb[1],
				CategoryLink: "/" + breadcrumb[1],
				Price:        productPrice,
				CurrentPrice: productCurrentPrice,
				Discount:     discount,
				Platform:     "setupgame",
			})
		})
	})
	d.Visit(url)
	d.Wait()
}
