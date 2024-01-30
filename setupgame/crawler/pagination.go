package crawler

import (
	"fmt"
	"github.com/gocolly/colly"
)

func GetURLs() {
	url := "https://setupgame.ma/"
	c := colly.NewCollector()
	c.OnHTML("div.elementor-element-33f27182:nth-child(1) > div:nth-child(2) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(6)", func(element *colly.HTMLElement) {
		fmt.Println(element)
	})
	c.Visit(url)
}
