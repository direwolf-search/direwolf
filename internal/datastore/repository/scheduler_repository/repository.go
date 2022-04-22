package schedulerrepository

import (
	"direwolf/internal/datastore/repository/engine"
	domainscheduler "direwolf/internal/domain/repository/scheduler"
	pkghelpers "direwolf/internal/pkg/helpers"
)

var (
	ErrFakeEngine       = "this is fake error from NewSchedulerRepository()"
	ErrRepositoryCreate = "error of create new scheduler repository"
)

func NewSchedulerRepository(dbEngine engine.DBEngine) (domainscheduler.Repository, error) {
	if dbEngine.GetEngineType() == engine.SQL.String() {
		if sqlEngine, ok := dbEngine.(domainscheduler.Repository); ok {
			return sqlEngine, nil
		}
	} else {
		return nil, pkghelpers.ErrorBuilder(ErrRepositoryCreate)
	}

	// other repository engine types can be injected here:
	// 		if dbEngine.GetEngineType() == engine.NoSQL.String() {
	//			if sqlEngine, ok := dbEngine.(domainscheduler.Repository); ok {
	//				return sqlEngine, nil
	//			} else {
	//				return nil, pkghelpers.ErrorBuilder(ErrRepositoryCreate)
	//			}
	//		}
	//

	return nil, pkghelpers.ErrorBuilder(ErrFakeEngine)
}
