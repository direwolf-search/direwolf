package repository

// Repository ...
type Repository interface {
	HostRepository
	LinkRepository
	CrawlerRepository
	EngineRepository
	SchedulerRepository
}
