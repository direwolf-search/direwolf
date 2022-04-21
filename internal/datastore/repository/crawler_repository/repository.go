package crawlerrepository

import (
	"direwolf/internal/datastore/repository/engine"
	domaincrawler "direwolf/internal/domain/repository/crawler"
	pkghelpers "direwolf/internal/pkg/helpers"
)

var (
	ErrFakeCrawler      = "this is fake error from NewCrawlerRepository()"
	ErrRepositoryCreate = "error of create new crawler repository"
)

func NewCrawlerRepository(dbEngine engine.DBEngine) (domaincrawler.Repository, error) {
	if dbEngine.GetEngineType() == engine.SQL.String() {
		if sqlEngine, ok := dbEngine.(domaincrawler.Repository); ok {
			return sqlEngine, nil
		} else {
			return nil, pkghelpers.ErrorBuilder(ErrRepositoryCreate)
		}
	}

	// other repository engine types can be injected here:
	// 		if dbEngine.GetEngineType() == engine.NoSQL.String() {
	//			if sqlEngine, ok := dbEngine.(domaincrawler.Repository); ok {
	//				return sqlEngine, nil
	//			} else {
	//				return nil, pkghelpers.ErrorBuilder(ErrRepositoryCreate)
	//			}
	//		}
	//

	return nil, pkghelpers.ErrorBuilder(ErrFakeCrawler)
}
