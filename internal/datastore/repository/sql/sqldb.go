package sql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
	//"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/schema"
)

var (
	noRowsErrMessage                 = "sql: no rows in result set"
	establishingConnectionErrMessage = "error establishing connection to repo:"
	hostHandlingErrMessage           = "error of host selecting or updating:"
	linkHandlingErrMessage           = "error of link selecting or updating:"
	notImplementedErrMessage         = "not implemented yet"
)

// errorBuilder builds error from error message and additional fields.
// Error message must be a first element in argument list.
func errorBuilder(fields ...interface{}) error {
	var formatString = ""
	for i := 0; i < len(fields); i++ {
		formatString += " %v"
	}
	return fmt.Errorf(formatString, fields...)
}

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
	db *bun.DB
}

// sqlRepository ...
type sqlRepository struct {
	Direct *bun.DB
	Repo   *SQLRepository
}

// NewSqlRepository ...
// TODO:NewSqlRepository(dialect schema.Dialect, dsn string)
func NewSqlRepository(dialect schema.Dialect) *sqlRepository {
	dsn := os.Getenv("DW_DEFAULT_TOR_CRAWLER_DSN")
	db := bun.NewDB(connToDB(dsn), dialect)

	if err := db.Ping(); err != nil {
		log.Fatalln(errorBuilder(establishingConnectionErrMessage, err.Error()))
	}

	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return &sqlRepository{
		Direct: db,
		Repo:   &SQLRepository{db: db},
	}
}

//func (sr *SQLRepository) Insert(ctx context.Context, entity interface{}) error {
//	switch v := entity.(type) {
//	case *host.Host:
//		return sr.CreateHost(ctx, v)
//	case *link.Link:
//		return sr.CreateLink(ctx, v)
//	}
//
//	return nil
//}
//
//// CreateHost ...
//func (sr *SQLRepository) CreateHost(ctx context.Context, h *host.Host) error {
//	var (
//		hostModel = sr.ConvertHostToModel(h)
//		id        = int64(-1)
//	)
//	// check if host already exists in DB
//	if err := sr.db.NewSelect().Model(hostModel).Column("id").
//		Where("? = ?", bun.Ident("url"), hostModel.URL).Scan(ctx); err != nil {
//		if !strings.Contains(err.Error(), noRowsErrMessage) {
//			wrappedErr := errorBuilder(hostHandlingErrMessage, err.Error(), hostModel.URL)
//			log.Println(wrappedErr.Error())
//			return wrappedErr
//		} else {
//			id = 0
//		}
//	}
//	// if host is not exists, insert it
//	if id == 0 {
//		if res, err := sr.db.NewInsert().Model(hostModel).Exec(ctx); err != nil {
//			wrappedErr := errorBuilder(hostHandlingErrMessage, err.Error(), hostModel.URL)
//			log.Println(wrappedErr.Error())
//			return wrappedErr
//		} else {
//			id, _ = res.LastInsertId()
//		}
//	} else {
//		id = hostModel.ID
//	}
//
//	return nil
//}
//
//// UpdateHostByURL ...
//func (sr *SQLRepository) UpdateHostByURL(ctx context.Context, url string, fields map[string]interface{}) error {
//	res, err := sr.db.NewUpdate().Model(&fields).TableExpr("hosts").Where("? = ?", bun.Ident("url"), url).
//		Exec(ctx)
//
//	if err != nil {
//		wrappedErr := errorBuilder(hostHandlingErrMessage, err.Error(), url)
//		log.Println(wrappedErr.Error())
//		return wrappedErr
//	}
//
//	num, _ := res.RowsAffected()
//	if num != 1 {
//		return errorBuilder(hostHandlingErrMessage, url)
//	}
//
//	return nil
//}
//
//// UpdateHostByID ...
//func (sr *SQLRepository) UpdateHostByID(ctx context.Context, id int64, fields map[string]interface{}) error {
//	res, err := sr.db.NewUpdate().Model(&fields).TableExpr("hosts").Where("? = ?", bun.Ident("id"), id).
//		Exec(ctx)
//
//	if err != nil {
//		wrappedErr := errorBuilder(hostHandlingErrMessage, err.Error(), id)
//		log.Println(wrappedErr.Error())
//		return wrappedErr
//	}
//
//	num, _ := res.RowsAffected()
//	if num != 1 {
//		return errorBuilder(hostHandlingErrMessage, id)
//	}
//
//	return nil
//}
//
//// GetHostByID ...
//func (sr *SQLRepository) GetHostByID(ctx context.Context, id int64) (*host.Host, error) {
//	var (
//		hostEntity = &host.Host{}
//	)
//	err := sr.db.NewSelect().Model((*models.Host)(nil)).Where("? = ?", bun.Ident("id"), id).Scan(ctx, hostEntity)
//	if err != nil {
//		wrappedErr := errorBuilder(hostHandlingErrMessage, err.Error(), id)
//		log.Println(wrappedErr.Error())
//		return nil, wrappedErr
//	}
//
//	return hostEntity, nil
//}
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
//		wrappedErr := errorBuilder(hostHandlingErrMessage, err.Error())
//		log.Println(wrappedErr.Error())
//		return nil, err
//	}
//
//	return hosts, nil
//}
//
//// GetAllHosts ...
//func (sr *SQLRepository) GetAllHosts(ctx context.Context) ([]*host.Host, error) {
//	var hosts = make([]*host.Host, 0)
//	if err := sr.db.NewSelect().Model((*models.Host)(nil)).
//		ColumnExpr("*").OrderExpr("id ASC").Scan(ctx, &hosts); err != nil {
//		wrappedErr := errorBuilder(hostHandlingErrMessage, err.Error())
//		log.Println(wrappedErr.Error())
//		return nil, wrappedErr
//	}
//
//	return hosts, nil
//}
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
//func (sr *SQLRepository) CreateLink(ctx context.Context, l *link.Link) error {
//	var (
//		linkModel = sr.ConvertLinkToModel(l)
//		id        = int64(-1)
//	)
//
//	// check if link already exists in DB
//	if err := sr.db.NewSelect().Model(linkModel).Column("id").
//		Where("? = ?", bun.Ident("body"), linkModel.Body).
//		Where("? = ?", bun.Ident("from"), linkModel.From).
//		Scan(ctx); err != nil {
//		if !strings.Contains(err.Error(), noRowsErrMessage) {
//			wrappedErr := errorBuilder(hostHandlingErrMessage, err.Error(), linkModel.String())
//			log.Println(wrappedErr.Error())
//			return wrappedErr
//		} else {
//			id = 0
//		}
//	}
//
//	if id == 0 {
//		if _, err := sr.db.NewInsert().Model(linkModel).Exec(ctx); err != nil {
//			wrappedErr := errorBuilder(hostHandlingErrMessage, err.Error(), linkModel.String())
//			log.Println(wrappedErr.Error())
//			return wrappedErr
//		}
//	}
//
//	return nil
//}
//
//// UpdateLink ...
//func (sr *SQLRepository) UpdateLink(ctx context.Context, id int64, fields map[string]interface{}) error {
//	res, err := sr.db.NewUpdate().Model(&fields).
//		TableExpr("links").
//		Where("? = ?", bun.Ident("id"), id).Exec(ctx)
//
//	if err != nil {
//		wrappedErr := errorBuilder(linkHandlingErrMessage, err.Error(), id)
//		log.Println(wrappedErr.Error())
//		return wrappedErr
//	}
//
//	num, _ := res.RowsAffected()
//	if num != 1 {
//		return errorBuilder(linkHandlingErrMessage, id)
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
//		wrappedErr := errorBuilder(hostHandlingErrMessage, err.Error(), id)
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
//		wrappedErr := errorBuilder(hostHandlingErrMessage, err.Error())
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
//		wrappedErr := errorBuilder(linkHandlingErrMessage, err.Error())
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
//		wrappedErr := errorBuilder(linkHandlingErrMessage, err.Error())
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
