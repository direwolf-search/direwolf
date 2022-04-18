package main

import (
	colly "direwolf/internal/services/crawler/ce"
	parser "direwolf/internal/services/crawler/html_parser"
	"log"

	"github.com/uptrace/bun/dialect/mysqldialect"

	sqlRepo "direwolf/internal/datastore/repository/sql"
)

func main() {
	repo := sqlRepo.NewSqlRepository(mysqldialect.New())
	engine := colly.NewCollyEngine(true, parser.NewHTMLParser(), repo.Repo, 0)

	log.Println("Crawler started")
}
