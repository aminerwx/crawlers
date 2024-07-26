package main

import (
	"fmt"

	"github.com/aminerwx/crawlers/cmd/atlasgaming"
)

func main() {
	// ultrapc.Runner()
	// setupgame.Runner()
	cat := atlasgaming.GetCategories()
	fmt.Println(cat)
	atlasgaming.Crawl(cat[0])
}
