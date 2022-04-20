package enginerepository

import (
	"direwolf/internal/datastore/helpers"
	"direwolf/internal/datastore/repository/engine/sql"
	"direwolf/internal/domain"
	"direwolf/internal/domain/repository/crawler/engine"
	pkghelpers "direwolf/internal/pkg/helpers"
)

var ErrFake = "this is fake error from NewEngineRepository()"

func NewEngineRepository(repositoryEngineType, dsn string, logger domain.Logger) (engine.Repository, error) {
	if repositoryEngineType == "sql" {
		db, err := helpers.ConnToDB(dsn, logger)
		if err != nil {
			return nil, err
		}

		repo, err := sql.NewRepositorySQL(logger, db)
		if err != nil {
			return nil, err
		}

		return repo, nil
	}

	//if repositoryEngineType == "no_sql" {
	//	// TODO: not implemented
	//}

	return nil, pkghelpers.ErrorBuilder(ErrFake)
}
