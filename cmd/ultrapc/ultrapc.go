package ultrapc

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"

	"github.com/aminerwx/crawlers/cmd/ultrapc/crawler"
	"github.com/aminerwx/crawlers/helper"
	"github.com/jackc/pgx/v5"
)

func Runner() {
	helper.Maybe(godotenv.Load())
	m, err := migrate.New(
		"file://migration",
		os.Getenv("DATABASE_URL"))
	helper.Maybe(err)
	m.Up()
	// links := menu.GetMenu()
	// combined_links := strings.Join(links, "\n")
	// b := []byte(combined_links)
	// helper.Maybe(os.WriteFile("ultrapc_links.txt", b, 0660))
	articles, err := crawler.Crawler("ultrapc/ultrapc_links.txt")
	helper.Maybe(err)
	// fmt.Println(len(articles))
	InsertBulk(articles, "components", os.Getenv("DATABASE_URL"))
	// DUMP to database
	//	fmt.Println(item)
}

func InsertBulk(articles []crawler.Article, table string, uri string) {
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
			"current_price",
			"discount",
			"stock",
			"name",
			"url",
			"category",
			"category_url",
			"subcategory",
			"subcategory_url",
			"menu",
			"menu_url",
			"platform",
		},
		pgx.CopyFromSlice(len(articles), func(i int) ([]any, error) {
			return []any{
				articles[i].Price,
				articles[i].CurrentPrice,
				articles[i].Discount,
				articles[i].Stock,
				articles[i].Name,
				articles[i].URL,
				articles[i].CategoryName,
				articles[i].CategoryURL,
				articles[i].SubCategoryName,
				articles[i].SubCategoryURL,
				articles[i].MenuName,
				articles[i].MenuURL,
				articles[i].Platform,
			}, nil
		}),
	)
	if err != nil {
		log.Printf("QueryRow failed: %v\nStderr: %v\n", err, os.Stderr)
		os.Exit(1)
	}
	fmt.Println("Total inserted rows = ", copyCount)
}
