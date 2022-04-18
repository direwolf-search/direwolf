package main

import (
	"log"

	"github.com/uptrace/bun/dialect/mysqldialect"

	sqlRepo "direwolf/internal/datastore/repository/sql"
	colly "direwolf/internal/services/crawler/engine/ce"
	parser "direwolf/internal/services/crawler/html_parser"
)

func main() {
	repo := sqlRepo.NewSqlRepository(mysqldialect.New())
	engine := colly.NewCollyEngine(true, parser.NewHTMLParser(), repo.Repo, 0)

	log.Println("Crawler started")
}
