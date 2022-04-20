package sql1

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
	//"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/schema"

	"direwolf/internal/datastore/models"
	"direwolf/internal/domain"
	"direwolf/internal/domain/model/host"
	"direwolf/internal/domain/model/link"
	"direwolf/internal/pkg/helpers"
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

func connToDB(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	return db
}

type SQLRepository struct {
	logger domain.Logger
	db     *bun.DB
}

// sqlRepository ...
type sqlRepository struct {
	Direct *bun.DB
	Repo   *SQLRepository
}

// NewSqlRepository ...
// TODO:NewSqlRepository(dialect schema.Dialect, dsn string)
func NewSqlRepository(dialect schema.Dialect, logger domain.Logger) (*sqlRepository, error) {
	dsn := os.Getenv("DW_DEFAULT_TOR_CRAWLER_DSN")
	db := bun.NewDB(connToDB(dsn), dialect)

	if err := db.Ping(); err != nil {
		logger.Fatal(helpers.ErrorBuilder(establishingConnectionErrMessage, err.Error()), "") // TODO: msg arg
		return nil, err
	}

	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return &sqlRepository{
		Direct: db,
		Repo: &SQLRepository{
			db:     db,
			logger: logger,
		},
	}, nil
}

func (sr *SQLRepository) insertHost(ctx context.Context, h *host.Host) error {
	var (
		host = models.NewHostFromModel(h)
	)
	// check if host already exists in DB
	exist, err := sr.hostExists(ctx, h.URL)
	if err != nil {
		err := helpers.ErrorBuilder(errorOfCheckIfHostExists, err.Error())
		sr.logger.Error(err, host.URL)
	}
	// if host is not exists, insert it
	if !exist {
		if _, err := sr.db.NewInsert().Model(host).Exec(ctx); err != nil {
			err := helpers.ErrorBuilder(hostHandlingErrMessage, err.Error())
			sr.logger.Error(err, host.URL)
			return err
		}
	}

	return nil
}

func (sr *SQLRepository) linkExists(ctx context.Context, linkBody, linkFrom string) (bool, error) {
	return sr.db.NewSelect().
		Model((*models.Link)(nil)).
		Where("body = ?", linkBody).
		Where("from = ?", linkFrom).
		Exists(ctx)
}

func (sr *SQLRepository) hostExists(ctx context.Context, url string) (bool, error) {
	return sr.db.NewSelect().
		Model((*models.Link)(nil)).
		Where("url = ?", url).
		Exists(ctx)
}

func (sr *SQLRepository) insertLink(ctx context.Context, l *link.Link) error {
	var (
		link = models.NewLinkFromModel(l)
	)
	// check if link already exists in DB
	exist, err := sr.linkExists(ctx, l.Body, l.From)
	if err != nil {
		err := helpers.ErrorBuilder(errorOfCheckIfLinkExists, err.Error())
		sr.logger.Error(err, "", map[string]interface{}{"link.Body": l.Body, "link.From": l.From})
		return err
	}

	if !exist {
		if _, err := sr.db.NewInsert().Model(link).Exec(ctx); err != nil {
			err := helpers.ErrorBuilder(errLinkInsert, err.Error())
			sr.logger.Error(err, "", map[string]interface{}{"link.Body": l.Body, "link.From": l.From})
			return err
		}
	}

	return nil
}

func (sr *SQLRepository) updateHostByURL(ctx context.Context, url string, fields map[string]interface{}) error {
	res, err := sr.db.NewUpdate().
		Model(&fields).
		TableExpr("hosts").
		Where("? = ?", bun.Ident("url"), url).
		Exec(ctx)

	if err != nil {
		err := helpers.ErrorBuilder(hostHandlingErrMessage, err.Error())
		sr.logger.Error(err, url)
		return err
	}

	num, _ := res.RowsAffected()
	if num != 1 {
		return helpers.ErrorBuilder(hostHandlingErrMessage, url)
	}

	return nil
}

func (sr *SQLRepository) getAllHosts(ctx context.Context) ([]*host.Host, error) {
	var hosts = make([]*host.Host, 0)
	if err := sr.db.NewSelect().
		Model((*models.Host)(nil)).
		ColumnExpr("*").
		OrderExpr("id ASC").Scan(ctx, &hosts); err != nil {
		err = helpers.ErrorBuilder(hostHandlingErrMessage, err.Error())
		sr.logger.Error(err, "")
		return nil, err
	}

	return hosts, nil
}

func (sr *SQLRepository) Insert(ctx context.Context, entity map[string]interface{}) error {
	if h, err := host.FromMap(entity); err != nil {
		err = helpers.ErrorBuilder(errorOfConvertMapToHost, err.Error())
		sr.logger.Error(err, "")
		return err
	} else {
		err = sr.insertHost(ctx, h)
		if err != nil {
			err = helpers.ErrorBuilder(errorInsertHost, err.Error())
			sr.logger.Error(err, "")
			return err
		}
	}

	if l, err := link.FromMap(entity); err != nil {
		err = helpers.ErrorBuilder(errorOfConvertMapToLink, err.Error())
		sr.logger.Error(err, "")
		return err
	} else {
		err = sr.insertLink(ctx, l)
		if err != nil {
			err = helpers.ErrorBuilder(errorInsertLink, err.Error())
			sr.logger.Error(err, "")
			return err
		}
	}

	return nil
}

func (sr *SQLRepository) Updated(ctx context.Context, url, md5hash string) (bool, error) {
	var (
		hash string
	)

	if err := sr.db.NewSelect().Model((*models.Host)(nil)).ColumnExpr("md5hash").Where("url = ?", url).
		Scan(ctx, &hash); err != nil {
		err = helpers.ErrorBuilder(hostHandlingErrMessage, err.Error(), url)
		sr.logger.Error(err, "")
		return false, err
	}

	if hash == "" {
		err := helpers.ErrorBuilder("md5hash cannot be empty", url)
		sr.logger.Error(err, "")
		return false, err
	}

	return hash == md5hash, nil
}

func (sr *SQLRepository) Exists(ctx context.Context, url string) (bool, error) {
	return sr.hostExists(ctx, url)
}

func (sr *SQLRepository) GetHostByID(ctx context.Context, id int64) (*host.Host, error) {
	var (
		hostEntity = &host.Host{}
	)
	err := sr.db.NewSelect().Model((*models.Host)(nil)).Where("? = ?", bun.Ident("id"), id).Scan(ctx, hostEntity)
	if err != nil {
		wrappedErr := helpers.ErrorBuilder(hostHandlingErrMessage, err.Error(), id)
		log.Println(wrappedErr.Error())
		return nil, wrappedErr
	}

	return hostEntity, nil
}

//// UpdateHostByURL ...

//func (sr *SQLRepository) Insert(ctx context.Context, entity interface{}) error {
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
//func (sr *SQLRepository) UpdateHostByID(ctx context.Context, id int64, fields map[string]interface{}) error {
//	res, err := sr.db.NewUpdate().Model(&fields).TableExpr("hosts").Where("? = ?", bun.Ident("id"), id).
//		Exec(ctx)
//
//	if err != nil {
//		wrappedErr := helpers.ErrorBuilder(hostHandlingErrMessage, err.Error(), id)
//		log.Println(wrappedErr.Error())
//		return wrappedErr
//	}
//
//	num, _ := res.RowsAffected()
//	if num != 1 {
//		return helpers.ErrorBuilder(hostHandlingErrMessage, id)
//	}
//
//	return nil
//}
//
//// GetHostByID ...

//
//// GetHostByFields ...
//func (sr *SQLRepository) GetHostByFields(ctx context.Context, fields map[string]interface{}) ([]*host.Host, error) {
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
//		wrappedErr := helpers.ErrorBuilder(hostHandlingErrMessage, err.Error())
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
//func (sr *SQLRepository) DeleteHost(ctx context.Context, id int64) error {
//	var h = &models.Host{}
//	_, err := sr.db.NewDelete().Model(h).Where("? = ?", bun.Ident("id"), id).Exec(ctx)
//
//	return err
//}
//
//// CreateLink ...

//
//// UpdateLink ...
//func (sr *SQLRepository) UpdateLink(ctx context.Context, id int64, fields map[string]interface{}) error {
//	res, err := sr.db.NewUpdate().Model(&fields).
//		TableExpr("links").
//		Where("? = ?", bun.Ident("id"), id).Exec(ctx)
//
//	if err != nil {
//		wrappedErr := helpers.ErrorBuilder(linkHandlingErrMessage, err.Error(), id)
//		log.Println(wrappedErr.Error())
//		return wrappedErr
//	}
//
//	num, _ := res.RowsAffected()
//	if num != 1 {
//		return helpers.ErrorBuilder(linkHandlingErrMessage, id)
//	}
//
//	return nil
//}
//
//// GetLinkByID ...
//func (sr *SQLRepository) GetLinkByID(ctx context.Context, id int64) (*link.Link, error) {
//	var (
//		linkEntity = &link.Link{}
//	)
//	err := sr.db.NewSelect().Model((*models.Link)(nil)).Where("? = ?", bun.Ident("id"), id).Scan(ctx, linkEntity)
//	if err != nil {
//		wrappedErr := helpers.ErrorBuilder(hostHandlingErrMessage, err.Error(), id)
//		log.Println(wrappedErr.Error())
//		return nil, wrappedErr
//	}
//
//	return linkEntity, nil
//}
//
//// GetLinkByFields ...
//func (sr *SQLRepository) GetLinkByFields(ctx context.Context, fields map[string]interface{}) ([]*link.Link, error) {
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
//		wrappedErr := helpers.ErrorBuilder(hostHandlingErrMessage, err.Error())
//		log.Println(wrappedErr.Error())
//		return nil, wrappedErr
//	}
//
//	return links, nil
//}
//
//// GetLinksByHost ...
//func (sr *SQLRepository) GetLinksByHost(ctx context.Context, id int64) ([]*link.Link, error) {
//	var links = make([]*link.Link, 0)
//	if err := sr.db.NewSelect().Model((*models.Link)(nil)).ColumnExpr("*").TableExpr("links AS l").
//		Join("LEFT JOIN hosts AS h").JoinOn("? = ?", bun.Ident("h.id"), bun.Ident("l.from_id")).
//		Where("?=?", bun.Ident("h.id"), id).Scan(ctx, &links); err != nil {
//		wrappedErr := helpers.ErrorBuilder(linkHandlingErrMessage, err.Error())
//		log.Println(wrappedErr.Error())
//		return nil, wrappedErr
//	}
//
//	return links, nil
//}
//
//// GetAllLinks ...
//func (sr *SQLRepository) GetAllLinks(ctx context.Context) ([]*link.Link, error) {
//	var links = make([]*link.Link, 0)
//	if err := sr.db.NewSelect().Model((*models.Link)(nil)).
//		ColumnExpr("*").OrderExpr("id ASC").Scan(ctx, &links); err != nil {
//		wrappedErr := helpers.ErrorBuilder(linkHandlingErrMessage, err.Error())
//		log.Println(wrappedErr.Error())
//		return nil, wrappedErr
//	}
//
//	return links, nil
//}
//
//// DeleteLink ...
//func (sr *SQLRepository) DeleteLink(ctx context.Context, id int64) error {
//	var l = &models.Link{}
//	_, err := sr.db.NewDelete().Model(l).Where("? = ?", bun.Ident("id"), id).Exec(ctx)
//
//	return err
//}
