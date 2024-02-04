package crawler

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aminerwx/crawlers/helper"
	"github.com/gocolly/colly"
)

type Article struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	ImageURL        string `json:"img_url"`
	CategoryName    string `json:"category"`
	CategoryURL     string `json:"category_url"`
	SubCategoryName string `json:"subcategory"`
	SubCategoryURL  string `json:"subcategory_url"`
	MenuName        string `json:"menu"`
	MenuURL         string `json:"menu_url"`
	Price           int    `json:"price"`
	Stock           int    `json:"stock"`
}

type ArticleSelector struct {
	Name     string `json:"name"`
	Price    string `json:"price"`
	URL      string `json:"url"`
	ImageURL string `json:"image_url"`
	Stock    string `json:"stock"`
}

func Crawler(linksFilePath string) ([]Article, error) {
	// Read links.txt
	data, err := os.ReadFile(linksFilePath)
	helper.Maybe(err)
	// fmt.Print(string(data))
	links := strings.Split(string(data), "\n")
	links = links[:len(links)-1]
	// fmt.Print(links)
	var productData []Article
	for i, url := range links {
		if i == 1 {
			break
		}

		pages, err := GetProductPages(url)
		fmt.Println("total pages = ", pages)
		helper.Maybe(err)
		for i := 1; i <= pages; i++ {
			var currentUrl strings.Builder
			currentUrl.WriteString(url)
			currentUrl.WriteString("?page=")
			currentUrl.WriteString(strconv.Itoa(i))
			articles, err := GetProduct(currentUrl.String())
			helper.Maybe(err)
			productData = append(productData, articles...)
		}
	}
	return productData, nil
}

func GetProduct(url string) ([]Article, error) {
	var breadcrumb [][]string
	var products []Article
	c := colly.NewCollector(
		colly.UserAgent(
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.79 Safari/537.36"),
		colly.MaxDepth(2),
		colly.Async(true),
	)

	helper.Maybe(c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2, RandomDelay: 4 * time.Second, Delay: 4}))

	c.OnError(func(_ *colly.Response, err error) {
		log.Println(err.Error())
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("--------------------- %s ---------------------\n", r.URL)
	})

	c.OnHTML("ol.list-unstyled", func(navBar *colly.HTMLElement) {
		navBar.ForEach("li.breadcrumb-item:not(:first-child)", func(_ int, child *colly.HTMLElement) {
			bread := []string{
				child.ChildText("a > span"),
				strings.ReplaceAll(child.ChildAttr("a", "href"), "https://www.ultrapc.ma", ""),
			}
			breadcrumb = append(breadcrumb, bread)
		})

		// fmt.Println(breadcrumb)
	})

	c.OnHTML(".products", func(h *colly.HTMLElement) {
		articleSelector := ArticleSelector{
			Name:     "div.product-block > div.product-left-block > h3.product-title > a",
			Price:    "div:nth-child(1) > div.product-right-block > div:nth-child(1) > span:nth-child(2)",
			URL:      "div:nth-child(1) > div:nth-child(1) > div:nth-child(2) > a:nth-child(1)",
			ImageURL: "div:nth-child(1) > div:nth-child(1) > div:nth-child(2) > a:nth-child(1) > img:nth-child(1)",
			Stock:    "div:nth-child(1) > div:nth-child(2) > div:nth-child(2) > strong:nth-child(2)",
		}
		pattern := regexp.MustCompile(`\d+`)
		h.ForEach("article.product-miniature", func(_ int, el *colly.HTMLElement) {
			article := Article{}
			price, _ := strconv.Atoi(pattern.FindString(helper.SpaceMap(el.ChildText(articleSelector.Price))))
			article.Name = el.ChildText(articleSelector.Name)
			article.Price = price
			article.URL = strings.ReplaceAll(el.ChildAttr(articleSelector.URL, "href"), "https://www.ultrapc.ma", "")
			article.ImageURL = strings.ReplaceAll(el.ChildAttr(articleSelector.ImageURL, "src"), "https://www.ultrapc.ma", "")
			stockText := el.DOM.Find(articleSelector.Stock).Text()
			if len(stockText) != 0 {
				stock, _ := strconv.Atoi(pattern.FindString(stockText))
				article.Stock = stock
			} else {
				article.Stock = 0
			}

			if len(breadcrumb) == 2 {
				article.MenuName = breadcrumb[0][0]
				article.MenuURL = breadcrumb[0][1]
				article.CategoryName = breadcrumb[1][0]
				article.CategoryURL = breadcrumb[1][1]
			} else {
				article.MenuName = breadcrumb[0][0]
				article.MenuURL = breadcrumb[0][1]
				article.CategoryName = breadcrumb[1][0]
				article.CategoryURL = breadcrumb[1][1]
				article.SubCategoryName = breadcrumb[2][0]
				article.SubCategoryURL = breadcrumb[2][1]
			}
			products = append(products, article)
		})
		// fmt.Printf("\n-> %s\n-> %s\n-> %s\n-> %s\n\n", name, price, url, image)
		breadcrumb = nil
	})
	err := c.Visit(url)
	if err != nil {
		log.Println(err)
		return []Article{}, nil
	}
	c.Wait()

	return products, nil
}

func GetProductPages(url string) (int, error) {
	pages := 0
	hasPages := false

	c := colly.NewCollector(
		colly.UserAgent(
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.79 Safari/537.36"),
		colly.MaxDepth(2),
		colly.Async(true),
	)
	helper.Maybe(c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2}))

	c.OnError(func(_ *colly.Response, err error) {
		log.Println(err.Error())
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("--------------------- %s ---------------------\n", r.URL)
	})

	c.OnHTML(".pagination", func(el *colly.HTMLElement) {
		el.ForEach("li.page-item", func(_ int, _ *colly.HTMLElement) {
			if !hasPages {
				pages++
			}
		})
		if pages > 1 {
			pages--
		}
		hasPages = true
	})

	helper.Maybe(c.Visit(url))
	c.Wait()
	return pages, nil
}
