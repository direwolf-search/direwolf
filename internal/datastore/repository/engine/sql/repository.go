package sql

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"

	//"direwolf/internal/datastore/helpers"
	"direwolf/internal/datastore/models"
	"direwolf/internal/domain"
	"direwolf/internal/domain/model/host"
	"direwolf/internal/domain/model/link"
	pkghelpers "direwolf/internal/pkg/helpers"
)

var (
	noRowsErrMessage                 = "sql: no rows in result set"
	errorOfConvertMapToHost          = "error of convert map to host"
	errorOfConvertMapToLink          = "error of convert map to link"
	errorOfCheckIfHostExists         = "error of check if host exists"
	errorOfCheckIfLinkExists         = "error of check if link exists"
	errorInsertHost                  = "error of host insert"
	errorInsertLink                  = "error of link insert"
	establishingConnectionErrMessage = "error establishing connection to repo:"
	hostHandlingErrMessage           = "error of host selecting or updating:"
	linkHandlingErrMessage           = "error of link selecting or updating:"
	notImplementedErrMessage         = "not implemented yet"
	errLinkInsert                    = "error of link inserting"
)

//type RepositorySQL struct {
//	logger domain.Logger
//	db     *bun.DB
//}
//
//// RepositorySQL ...
//type RepositorySQL struct {
//	Direct *bun.DB
//	Repo   *RepositorySQL
//}
//
//// NewRepositorySQL ...
//// TODO:NewRepositorySQL(dialect schema.Dialect, dsn string)
//func NewRepositorySQL(dialect schema.Dialect, logger domain.Logger) (*RepositorySQL, error) {
//	dsn := os.Getenv("DW_DEFAULT_TOR_CRAWLER_DSN")
//	db := bun.NewDB(connToDB(dsn), dialect)
//
//	if err := db.Ping(); err != nil {
//		logger.Fatal(pkghelpers.ErrorBuilder(establishingConnectionErrMessage, err.Error()), "") // TODO: msg arg
//		return nil, err
//	}
//
//	// Print all queries to stdout.
//	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
//
//	return &RepositorySQL{
//		Direct: db,
//		Repo: &RepositorySQL{
//			db:     db,
//			logger: logger,
//		},
//	}, nil
//}

type RepositorySQL struct {
	logger domain.Logger
	db     *bun.DB
}

func NewRepositorySQL(logger domain.Logger, sqlDB *sql.DB) (*RepositorySQL, error) {
	db := bun.NewDB(sqlDB, mysqldialect.New())

	if err := db.Ping(); err != nil {
		logger.Fatal(pkghelpers.ErrorBuilder("error establishing connection to bun.DB:", err.Error()), "") // TODO: msg arg
		return nil, err
	}

	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return &RepositorySQL{
		logger: logger,
		db:     db,
	}, nil
}

func (sr *RepositorySQL) InsertHost(ctx context.Context, h *host.Host) error {
	var (
		host = models.NewHostFromModel(h)
	)
	// check if host already exists in DB
	exist, err := sr.HostExists(ctx, h.URL)
	if err != nil {
		err := pkghelpers.ErrorBuilder(errorOfCheckIfHostExists, err.Error())
		sr.logger.Error(err, host.URL)
	}
	// if host is not exists, insert it
	if !exist {
		if _, err := sr.db.NewInsert().Model(host).Exec(ctx); err != nil {
			err := pkghelpers.ErrorBuilder(hostHandlingErrMessage, err.Error())
			sr.logger.Error(err, host.URL)
			return err
		}
	}

	return nil
}

func (sr *RepositorySQL) LinkExists(ctx context.Context, linkBody, linkFrom string) (bool, error) {
	return sr.db.NewSelect().
		Model((*models.Link)(nil)).
		Where("body = ?", linkBody).
		Where("from = ?", linkFrom).
		Exists(ctx)
}

func (sr *RepositorySQL) HostExists(ctx context.Context, url string) (bool, error) {
	return sr.db.NewSelect().
		Model((*models.Link)(nil)).
		Where("url = ?", url).
		Exists(ctx)
}

func (sr *RepositorySQL) InsertLink(ctx context.Context, l *link.Link) error {
	var (
		link = models.NewLinkFromModel(l)
	)
	// check if link already exists in DB
	exist, err := sr.LinkExists(ctx, l.Body, l.From)
	if err != nil {
		err := pkghelpers.ErrorBuilder(errorOfCheckIfLinkExists, err.Error())
		sr.logger.Error(err, "", map[string]interface{}{"link.Body": l.Body, "link.From": l.From})
		return err
	}

	if !exist {
		if _, err := sr.db.NewInsert().Model(link).Exec(ctx); err != nil {
			err := pkghelpers.ErrorBuilder(errLinkInsert, err.Error())
			sr.logger.Error(err, "", map[string]interface{}{"link.Body": l.Body, "link.From": l.From})
			return err
		}
	}

	return nil
}

func (sr *RepositorySQL) UpdateHostByURL(ctx context.Context, url string, fields map[string]interface{}) error {
	res, err := sr.db.NewUpdate().
		Model(&fields).
		TableExpr("hosts").
		Where("? = ?", bun.Ident("url"), url).
		Exec(ctx)

	if err != nil {
		err := pkghelpers.ErrorBuilder(hostHandlingErrMessage, err.Error())
		sr.logger.Error(err, url)
		return err
	}

	num, _ := res.RowsAffected()
	if num != 1 {
		return pkghelpers.ErrorBuilder(hostHandlingErrMessage, url)
	}

	return nil
}

func (sr *RepositorySQL) GetAllHosts(ctx context.Context) ([]*host.Host, error) {
	var hosts = make([]*host.Host, 0)
	if err := sr.db.NewSelect().
		Model((*models.Host)(nil)).
		ColumnExpr("*").
		OrderExpr("id ASC").Scan(ctx, &hosts); err != nil {
		err = pkghelpers.ErrorBuilder(hostHandlingErrMessage, err.Error())
		sr.logger.Error(err, "")
		return nil, err
	}

	return hosts, nil
}

func (sr *RepositorySQL) InsertHostOrLink(ctx context.Context, entity map[string]interface{}) error {
	if h, err := host.FromMap(entity); err != nil {
		err = pkghelpers.ErrorBuilder(errorOfConvertMapToHost, err.Error())
		sr.logger.Error(err, "")
		return err
	} else {
		err = sr.InsertHost(ctx, h)
		if err != nil {
			err = pkghelpers.ErrorBuilder(errorInsertHost, err.Error())
			sr.logger.Error(err, "")
			return err
		}
	}

	if l, err := link.FromMap(entity); err != nil {
		err = pkghelpers.ErrorBuilder(errorOfConvertMapToLink, err.Error())
		sr.logger.Error(err, "")
		return err
	} else {
		err = sr.InsertLink(ctx, l)
		if err != nil {
			err = pkghelpers.ErrorBuilder(errorInsertLink, err.Error())
			sr.logger.Error(err, "")
			return err
		}
	}

	return nil
}

// *RepositorySQL implements engine.Repository interface

func (sr *RepositorySQL) Insert(ctx context.Context, entity map[string]interface{}) error {
	return sr.InsertHostOrLink(ctx, entity)
}

func (sr *RepositorySQL) Updated(ctx context.Context, url, md5hash string) (bool, error) {
	return sr.HostUpdated(ctx, url, md5hash)
}

func (sr *RepositorySQL) Exists(ctx context.Context, url string) (bool, error) {
	return sr.HostExists(ctx, url)
}

func (sr *RepositorySQL) Update(ctx context.Context, entity map[string]interface{}) error {
	return sr.UpdateHost(ctx, entity)
}

func (sr *RepositorySQL) HostUpdated(ctx context.Context, url, md5hash string) (bool, error) {
	var (
		hash string
	)

	if err := sr.db.NewSelect().Model((*models.Host)(nil)).ColumnExpr("md5hash").Where("url = ?", url).
		Scan(ctx, &hash); err != nil {
		err = pkghelpers.ErrorBuilder(hostHandlingErrMessage, err.Error(), url)
		sr.logger.Error(err, "")
		return false, err
	}

	if hash == "" {
		err := pkghelpers.ErrorBuilder("md5hash cannot be empty", url)
		sr.logger.Error(err, "")
		return false, err
	}

	return hash == md5hash, nil
}

func (sr *RepositorySQL) UpdateHost(ctx context.Context, entity map[string]interface{}) error {
	if _, err := host.FromMap(entity); err != nil {
		err = pkghelpers.ErrorBuilder(errorOfConvertMapToHost, err.Error())
		sr.logger.Error(err, "")
		return err
	}

	res, err := sr.db.NewUpdate().
		Model(&entity).
		TableExpr("hosts").
		OmitZero().
		Where("url = ?", entity["url"]).
		Exec(ctx)

	if err != nil {
		wrappedErr := pkghelpers.ErrorBuilder(hostHandlingErrMessage, err.Error(), entity["url"])
		sr.logger.Error(err, "")
		return wrappedErr
	}

	num, _ := res.RowsAffected()
	if num != 1 {
		err = pkghelpers.ErrorBuilder(hostHandlingErrMessage, entity["url"])
		sr.logger.Error(err, "")
		return err
	}

	return nil
}

func (sr *RepositorySQL) GetHostByID(ctx context.Context, id int64) (*host.Host, error) {
	var (
		hostEntity = &host.Host{}
	)
	err := sr.db.NewSelect().Model((*models.Host)(nil)).Where("? = ?", bun.Ident("id"), id).Scan(ctx, hostEntity)
	if err != nil {
		wrappedErr := pkghelpers.ErrorBuilder(hostHandlingErrMessage, err.Error(), id)
		log.Println(wrappedErr.Error())
		return nil, wrappedErr
	}

	return hostEntity, nil
}

//// UpdateHostByURL ...

//func (sr *RepositorySQL) Insert(ctx context.Context, entity interface{}) error {
//	switch v := entity.(type) {
//	case *host.Host:
//		return sr.CreateHost(ctx, v)
//	case *link.Link:x
//		return sr.CreateLink(ctx, v)
//	}
//
//	return nil
//}
//
//// CreateHost ...

//
//// UpdateHostByID ...
//func (sr *RepositorySQL) UpdateHostByID(ctx context.Context, id int64, fields map[string]interface{}) error {
//	res, err := sr.db.NewUpdate().Model(&fields).TableExpr("hosts").Where("? = ?", bun.Ident("id"), id).
//		Exec(ctx)
//
//	if err != nil {
//		wrappedErr := pkghelpers.ErrorBuilder(hostHandlingErrMessage, err.Error(), id)
//		log.Println(wrappedErr.Error())
//		return wrappedErr
//	}
//
//	num, _ := res.RowsAffected()
//	if num != 1 {
//		return pkghelpers.ErrorBuilder(hostHandlingErrMessage, id)
//	}
//
//	return nil
//}
//
//// GetHostByID ...

//
//// GetHostByFields ...
//func (sr *RepositorySQL) GetHostByFields(ctx context.Context, fields map[string]interface{}) ([]*host.Host, error) {
//	var (
//		hosts        = make([]*host.Host, 0)
//		whereBuilder = func(
//			ctx context.Context,
//			q *bun.SelectQuery,
//			fields map[string]interface{},
//			hosts []*host.Host,
//		) error {
//			for fieldName, fieldValue := range fields {
//				q = q.Where("? = ?", bun.Ident(fieldName), fieldValue)
//			}
//			return q.Scan(ctx, &hosts)
//		}
//	)
//	q := sr.db.NewSelect().Model((*models.Host)(nil)).ColumnExpr("*")
//	if err := whereBuilder(ctx, q, fields, hosts); err != nil {
//		wrappedErr := pkghelpers.ErrorBuilder(hostHandlingErrMessage, err.Error())
//		log.Println(wrappedErr.Error())
//		return nil, err
//	}
//
//	return hosts, nil
//}
//
//// GetAllHosts ...

//
//// DeleteHost ...
//func (sr *RepositorySQL) DeleteHost(ctx context.Context, id int64) error {
//	var h = &models.Host{}
//	_, err := sr.db.NewDelete().Model(h).Where("? = ?", bun.Ident("id"), id).Exec(ctx)
//
//	return err
//}
//
//// CreateLink ...

//
//// UpdateLink ...
//func (sr *RepositorySQL) UpdateLink(ctx context.Context, id int64, fields map[string]interface{}) error {
//	res, err := sr.db.NewUpdate().Model(&fields).
//		TableExpr("links").
//		Where("? = ?", bun.Ident("id"), id).Exec(ctx)
//
//	if err != nil {
//		wrappedErr := pkghelpers.ErrorBuilder(linkHandlingErrMessage, err.Error(), id)
//		log.Println(wrappedErr.Error())
//		return wrappedErr
//	}
//
//	num, _ := res.RowsAffected()
//	if num != 1 {
//		return pkghelpers.ErrorBuilder(linkHandlingErrMessage, id)
//	}
//
//	return nil
//}
//
//// GetLinkByID ...
//func (sr *RepositorySQL) GetLinkByID(ctx context.Context, id int64) (*link.Link, error) {
//	var (
//		linkEntity = &link.Link{}
//	)
//	err := sr.db.NewSelect().Model((*models.Link)(nil)).Where("? = ?", bun.Ident("id"), id).Scan(ctx, linkEntity)
//	if err != nil {
//		wrappedErr := pkghelpers.ErrorBuilder(hostHandlingErrMessage, err.Error(), id)
//		log.Println(wrappedErr.Error())
//		return nil, wrappedErr
//	}
//
//	return linkEntity, nil
//}
//
//// GetLinkByFields ...
//func (sr *RepositorySQL) GetLinkByFields(ctx context.Context, fields map[string]interface{}) ([]*link.Link, error) {
//	var (
//		links        = make([]*link.Link, 0)
//		whereBuilder = func(
//			ctx context.Context,
//			q *bun.SelectQuery,
//			fields map[string]interface{},
//			links []*link.Link,
//		) error {
//			for fieldName, fieldValue := range fields {
//				q = q.Where("? = ?", bun.Ident(fieldName), fieldValue)
//			}
//			return q.Scan(ctx, &links)
//		}
//	)
//	q := sr.db.NewSelect().Model((*models.Link)(nil)).ColumnExpr("*")
//	if err := whereBuilder(ctx, q, fields, links); err != nil {
//		wrappedErr := pkghelpers.ErrorBuilder(hostHandlingErrMessage, err.Error())
//		log.Println(wrappedErr.Error())
//		return nil, wrappedErr
//	}
//
//	return links, nil
//}
//
//// GetLinksByHost ...
//func (sr *RepositorySQL) GetLinksByHost(ctx context.Context, id int64) ([]*link.Link, error) {
//	var links = make([]*link.Link, 0)
//	if err := sr.db.NewSelect().Model((*models.Link)(nil)).ColumnExpr("*").TableExpr("links AS l").
//		Join("LEFT JOIN hosts AS h").JoinOn("? = ?", bun.Ident("h.id"), bun.Ident("l.from_id")).
//		Where("?=?", bun.Ident("h.id"), id).Scan(ctx, &links); err != nil {
//		wrappedErr := pkghelpers.ErrorBuilder(linkHandlingErrMessage, err.Error())
//		log.Println(wrappedErr.Error())
//		return nil, wrappedErr
//	}
//
//	return links, nil
//}
//
//// GetAllLinks ...
//func (sr *RepositorySQL) GetAllLinks(ctx context.Context) ([]*link.Link, error) {
//	var links = make([]*link.Link, 0)
//	if err := sr.db.NewSelect().Model((*models.Link)(nil)).
//		ColumnExpr("*").OrderExpr("id ASC").Scan(ctx, &links); err != nil {
//		wrappedErr := pkghelpers.ErrorBuilder(linkHandlingErrMessage, err.Error())
//		log.Println(wrappedErr.Error())
//		return nil, wrappedErr
//	}
//
//	return links, nil
//}
//
//// DeleteLink ...
//func (sr *RepositorySQL) DeleteLink(ctx context.Context, id int64) error {
//	var l = &models.Link{}
//	_, err := sr.db.NewDelete().Model(l).Where("? = ?", bun.Ident("id"), id).Exec(ctx)
//
//	return err
//}
