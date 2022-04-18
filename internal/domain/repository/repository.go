package repository

// Repository ...
type Repository interface {
	CrawlerRepository
	EngineRepository
	SchedulerRepository
}
