package crawler

type Crawler interface {
	Crawl(links []string)
}
