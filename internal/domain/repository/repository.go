package repository

// Repository ...
type Repository interface {
	HostRepository
	LinkRepository
	CrawlerTaskPoolRepository
}
