package enginerepository

import (
	"direwolf/internal/datastore/repository/engine"
	domainengine "direwolf/internal/domain/repository/crawler/engine"
	pkghelpers "direwolf/internal/pkg/helpers"
)

var (
	ErrFakeEngine       = "this is fake error from NewEngineRepository()"
	ErrRepositoryCreate = "error of create new crawler repository"
)

func NewEngineRepository(dbEngine engine.DBEngine) (domainengine.Repository, error) {
	if dbEngine.GetEngineType() == engine.SQL.String() {
		if sqlEngine, ok := dbEngine.(domainengine.Repository); ok {
			return sqlEngine, nil
		} else {
			return nil, pkghelpers.ErrorBuilder(ErrRepositoryCreate)
		}
	}

	// other repository engine types can be injected here:
	// 		if dbEngine.GetEngineType() == engine.NoSQL.String() {
	//			if sqlEngine, ok := dbEngine.(domainengine.Repository); ok {
	//				return sqlEngine, nil
	//			} else {
	//				return nil, pkghelpers.ErrorBuilder(ErrRepositoryCreate)
	//			}
	//		}
	//

	return nil, pkghelpers.ErrorBuilder(ErrFakeEngine)
}
