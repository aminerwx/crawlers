package setupgame

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aminerwx/crawlers/helper"
	"github.com/aminerwx/crawlers/setupgame/crawler"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func Runner() {
	links := crawler.GetLinks()
	var products []crawler.Product
	for i, link := range links {
		fmt.Println("###", link)
		crawler.Crawler(link, &products)
		if i == 1 {
			break
		}
	}
	fmt.Println(products)
	// crawler.Crawler("https://setupgame.ma/categorie-produit/composants-gaming/boitier-gamer/")
}

func Init(products *[]crawler.Product) {
	_, err := migrate.New(
		"file://db/migrations",
		"postgres://postgres@localhost:5432/ultrapc?sslmode=disable")
	helper.Maybe(err)
	helper.Maybe(godotenv.Load())
	// m.Up()
	// m.Down()
	// links := menu.GetMenu()
	// combined_links := strings.Join(links, "\n")
	// b := []byte(combined_links)
	// helper.Maybe(os.WriteFile("ultrapc_links.txt", b, 0660))
	helper.Maybe(err)
	// fmt.Println(len(articles))
	InsertBulk(products, "components", os.Getenv("DATABASE_URL"))
	// DUMP to database
	//	fmt.Println(item)
}

func InsertBulk(products *[]crawler.Product, table string, uri string) {
	conn, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		log.Printf("Unable to connect to database: %v\nStderr: %v\n", err, os.Stderr)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	copyCount, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{table},
		[]string{
			"price",
			"stock",
			"name",
			"url",
			"img_url",
			"category",
			"category_url",
			"subcategory",
			"subcategory_url",
			"menu",
			"menu_url",
		},
		pgx.CopyFromSlice(len(*products), func(i int) ([]any, error) {
			return []any{
				(*products)[i].Price,
				//(*products)[i].Stock,
				(*products)[i].Name,
				//(*products)[i].URL,
				//(*products)[i].ImageURL,
				//(*products)[i].CategoryName,
				//(*products)[i].CategoryURL,
				//(*products)[i].SubCategoryName,
				//(*products)[i].SubCategoryURL,
				//(*products)[i].MenuName,
				//(*products)[i].MenuURL,
			}, nil
		}),
	)
	if err != nil {
		log.Printf("QueryRow failed: %v\nStderr: %v\n", err, os.Stderr)
		os.Exit(1)
	}
	fmt.Println("Total inserted rows = ", copyCount)
}
