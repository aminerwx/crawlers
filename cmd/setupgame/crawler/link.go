package crawler

import (
	"github.com/gocolly/colly"
)

func RemoveDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	var list []T
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func GetLinks() []string {
	c := colly.NewCollector()
	wantedMenu := []string{
		"#home-menu > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(6)",
		"#home-menu > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(7)",
		"#home-menu > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(8)",
		"#home-menu > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(9)",
	}
	var linkItems []string
	// var linkCache = make(map[string]bool)
	// menuSelector := wantedMenu[0]
	for _, menuSelector := range wantedMenu {
		GetCategories(c, menuSelector, &linkItems)
	}

	c.Visit("https://setupgame.ma/")

	return RemoveDuplicate(linkItems)
}

func GetCategories(c *colly.Collector, menuSelector string, linkItems *[]string) {
	c.OnHTML(menuSelector, func(element *colly.HTMLElement) {
		categorySelector := menuSelector + " div:nth-child(2) > div"
		element.ForEach(categorySelector, func(_ int, category *colly.HTMLElement) {
			childNodes := category.DOM.Find("div:nth-child(2)").Nodes
			// categoryLink := strings.ReplaceAll(category.DOM.Find("a:nth-child(1)").AttrOr("href", "-"), linkReplacer, "")
			if len(childNodes) == 2 {
				category.ForEach("div:nth-child(2) > div", func(_ int, child *colly.HTMLElement) {
					subCategorySelector := "a:nth-child(1)"
					subCategoryURL := child.DOM.Find(subCategorySelector).AttrOr("href", "-")
					//if !linkCache[subCategoryURL] {
					//	linkItems = append(linkItems, subCategoryURL)
					//	linkCache[subCategoryURL] = true
					//}
					*linkItems = append(*linkItems, subCategoryURL)
				})
			}
			if len(childNodes) == 0 {
				categoryLink := category.DOM.Find("a:nth-child(1)").AttrOr("href", "-")
				//if !linkCache[categoryLink] {
				//	linkItems = append(linkItems, categoryLink)
				//	linkCache[categoryLink] = true
				//}
				*linkItems = append(*linkItems, categoryLink)
			}
		})
	})
}
