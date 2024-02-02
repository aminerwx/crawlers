package helper

import (
	"log"
	"strings"
	"unicode"
)

func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func Maybe(err error) {
	if err != nil {
		log.Fatalf("Something is wrong : %s", err)
	}
}

func ParseLink(link string) []string {
	baseLink := "https://setupgame.ma/categorie-produit/"
	return strings.Split(strings.ReplaceAll(link, baseLink, ""), "/")
}
