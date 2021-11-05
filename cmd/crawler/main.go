package main

import (
	"log"

	"github.com/uptrace/bun/dialect/mysqldialect"

	sqlRepo "direwolf/internal/datastore/repository/sql"
	colly "direwolf/internal/pkg/crawler/ce"
	parser "direwolf/internal/pkg/crawler/html_parser"
)

func main() {
	repo := sqlRepo.NewSqlRepository(mysqldialect.New())
	engine := colly.NewCollyEngine(true, parser.NewHTMLParser(), repo.Repo, 0)

	log.Println("Crawler started")
}
