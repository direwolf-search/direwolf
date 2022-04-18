package repository

import (
	"direwolf/internal/domain/repository/crawler"
	"direwolf/internal/domain/repository/crawler/engine"
	"direwolf/internal/domain/repository/scheduler"
)

// Repository ...
type Repository interface {
	crawler.Repository
	engine.Repository
	scheduler.Repository
}
