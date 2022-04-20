package helpers

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"direwolf/internal/domain"
	pkghelpers "direwolf/internal/pkg/helpers"
)

var (
	errorDBConnection = "error of connection to DB"
	errorPingDB       = "error of ping DB"
)

func ConnToDB(dsn string, logger domain.Logger) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		err = pkghelpers.ErrorBuilder(errorDBConnection, err.Error())
		logger.Fatal(err, "")
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		err = pkghelpers.ErrorBuilder(errorPingDB, err.Error())
		logger.Fatal(err, "")
		return nil, err
	}

	return db, nil
}
