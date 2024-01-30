package crawler

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

func Crawler() {
	fmt.Println("calling Crawler()...")
}
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

// submenu			 categorySelector	categorysubmenusel
// "div:nth-child(2) > div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > a:nth-child(1)"
// "div:nth-child(2) > div:nth-child(1) > div:nth-child(2) > div:nth-child(2) > a:nth-child(1)"
// "div:nth-child(2) > div:nth-child(1) > a:nth-child(1) > span:nth-child(1) > span:nth-child(1)"
// div:nth-child(2) > div:nth-child(1)
// div:nth-child(2) > div:nth-child(2)

func GetLinks() []string {
	c := colly.NewCollector()
	wantedMenu := []string{
		"#home-menu > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(6)",
		"#home-menu > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(7)",
		"#home-menu > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(8)",
		"#home-menu > div:nth-child(1) > div:nth-child(1) > div:nth-child(1) > div:nth-child(9)",
	}
	var linkItems []string
	//var linkCache = make(map[string]bool)
	//menuSelector := wantedMenu[0]
	for _, menuSelector := range wantedMenu {
		GetData(c, menuSelector, &linkItems)
	}

	c.Visit("https://setupgame.ma/")

	return RemoveDuplicate(linkItems)
}
func GetData(c *colly.Collector, menuSelector string, linkItems *[]string) {
	c.OnHTML(menuSelector, func(element *colly.HTMLElement) {
		categorySelector := menuSelector + " div:nth-child(2) > div"
		element.ForEach(categorySelector, func(i int, category *colly.HTMLElement) {
			childNodes := category.DOM.Find("div:nth-child(2)").Nodes
			//categoryLink := strings.ReplaceAll(category.DOM.Find("a:nth-child(1)").AttrOr("href", "-"), linkReplacer, "")
			if len(childNodes) == 2 {
				category.ForEach("div:nth-child(2) > div", func(i int, child *colly.HTMLElement) {
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
func getData(c *colly.Collector, linkItems *[]string) {
	linkReplacer := "https://setupgame.ma/categorie-produit"
	menuSelector := "#home-menu > div > div > div > div:nth-child(6)"
	c.OnHTML(menuSelector, func(element *colly.HTMLElement) {
		// composants
		//subMenuSelector := menuSelector + " > div:nth-child(2)"

		// subcomposants cpu / mobo / ram ....
		// div > div.menu-item.menu-item-type-taxonomy.menu-item-object-product_cat.menu-item-has-children.jet-custom-nav__item.jet-custom-nav__item-66463.purple-border
		categorySelector := "div:nth-child(2) > div"
		element.ForEach(categorySelector, func(i int, category *colly.HTMLElement) {
			childNodes := category.DOM.Find("div").Nodes
			categoryLink := strings.ReplaceAll(category.DOM.Find("a").AttrOr("href", "-"), linkReplacer, "")
			if len(childNodes) > 0 {
				subCategorySelector := "div > div"
				category.ForEach(subCategorySelector, func(i int, child *colly.HTMLElement) {
					subCategoryLinkSelector := "a:nth-child(1)"
					subCategoryLink := strings.ReplaceAll(child.DOM.Find(subCategoryLinkSelector).AttrOr("href", "-"), linkReplacer, "")

					//names := strings.Split(child.Text, ":")
					////fmt.Println(names)
					//if len(names) == 2 {
					//	linkItem.CategoryName = names[0]
					//	linkItem.SubCategoryName = names[1]
					//}
					//if len(names) == 1 {
					//	linkItem.SubCategoryName = names[0]
					//}
					*linkItems = append(*linkItems, subCategoryLink)
				})
			} else {
				*linkItems = append(*linkItems, categoryLink)
			}
		})
		//fmt.Println(linkItems)
	})
}
