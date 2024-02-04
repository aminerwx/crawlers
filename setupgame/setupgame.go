package setupgame

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aminerwx/crawlers/helper"
	"github.com/aminerwx/crawlers/setupgame/crawler"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func Runner() {
	links := crawler.GetLinks()
	var products []crawler.Product
	for i, link := range links {
		fmt.Println("###", link)
		crawler.Crawler(link, &products)
		if i == 0 {
			break
		}
	}
	Init(&products)
	// fmt.Println(products)
}

func Init(products *[]crawler.Product) {
	m, err := migrate.New(
		"file://migration",
		"postgres://postgres@localhost:5432/crawlers?sslmode=disable")
	helper.Maybe(err)
	helper.Maybe(godotenv.Load())
	m.Down()
	m.Up()
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
			"name",
			"url",
			"menu",
			"menu_url",
			"category",
			"category_url",
			"platform",
			"price",
			"current_price",
			"discount",
			"stock",
		},
		/*	Name         string
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
		*/
		pgx.CopyFromSlice(len(*products), func(i int) ([]any, error) {
			return []any{
				(*products)[i].Name,
				(*products)[i].Link,
				(*products)[i].Menu,
				(*products)[i].MenuLink,
				(*products)[i].Category,
				(*products)[i].CategoryLink,
				(*products)[i].Platform,
				(*products)[i].Price,
				(*products)[i].CurrentPrice,
				(*products)[i].Discount,
				(*products)[i].Stock,
			}, nil
		}),
	)
	if err != nil {
		log.Printf("QueryRow failed: %v\nStderr: %v\n", err, os.Stderr)
		os.Exit(1)
	}
	fmt.Println("Total inserted rows = ", copyCount)
}
