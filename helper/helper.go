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

//func CategoryID(str string) []string {
//	categories := []string{
//		"Claviers",             // 1
//		"Bras & Pied",          // 2
//		"Cables",               // 3
//		"Souris",               // 4
//		"Kits claviers/souris", // 5
//		"Casques",              // 6
//		"Microphones",          // 7
//		"Tapis de souris",      // 8
//		"Volants pour PC",      // 9
//		"Réalité Virtuelle",
//		"Enceinte PC",
//		"Webcams",
//		"Chaise&Bureau",
//		"Bureau gamer",
//		"Onduleurs",
//		"Réseau",
//		"Stockage externe",
//		"Périphériques de jeu",
//		"Tablettes",
//		"Imprimantes/scanner",
//		"Connectique",
//		"Disques durs et SSD",
//		"Processeurs",
//		"Cartes mères",
//		"Refroidissement",
//		"Cartes graphiques",
//		"Mémoire vive PC",
//		"Alimentations PC",
//		"Boitiers PC",
//		"Cartes son",
//	}
//
//	//for i, cat := range categories {
//	//	if str == cat {
//	//		return i
//	//	}
//	//}
//	return categories
//}
