package setupgame

import (
	"fmt"
	"github.com/aminerwx/crawlers/setupgame/crawler"
)

func SetupRunner() {
	//crawler.GetLinks()
	links := crawler.GetLinks()
	for _, link := range links {
		fmt.Println("###", link)
	}
	crawler.Crawler()
}
