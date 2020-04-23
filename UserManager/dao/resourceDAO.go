package dao

import (
	"DBproject1/core"
	"DBproject1/model"
	"context"

	"github.com/linxGnu/mssqlx"
)

const (
	sqlResourceSave         = "INSERT INTO resource(name,description) VALUES (?,?)"
	sqlResourceDelete       = "DELETE FROM resource WHERE name=?"
	sqlResourceSelectAll    = "SELECT *FROM resource"
	sqlResourceSelectByView = "SELECT *FROM resource ORDER BY id DESC LIMIT ?, ?"
)

type IResourceDAO interface {
	SaveResource(ctx context.Context, db *mssqlx.DBs, name, desc string) (err error)
	DeleteResource(ctx context.Context, db *mssqlx.DBs, name string) (err error)
	SelectAllResource(ctx context.Context, db *mssqlx.DBs) (result []*model.Resource, err error)
	SelectResourceByView(ctx context.Context, db *mssqlx.DBs, offset, limit int64) (result []*model.Resource, err error)
}

type resourceDAO struct{}

func (c *resourceDAO) SaveResource(ctx context.Context, db *mssqlx.DBs, name, desc string) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}
	_, err = db.ExecContext(ctx, sqlResourceSave, name, desc)
	return
}

func (c *resourceDAO) DeleteResource(ctx context.Context, db *mssqlx.DBs, name string) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}
	_, err = db.ExecContext(ctx, sqlResourceDelete, name)
	return
}

func (c *resourceDAO) SelectAllResource(ctx context.Context, db *mssqlx.DBs) (result []*model.Resource, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}
	err = db.SelectContext(ctx, &result, sqlResourceSelectAll)
	return
}

func (c *resourceDAO) SelectResourceByView(ctx context.Context, db *mssqlx.DBs, offset, limit int64) (result []*model.Resource, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlResourceSelectByView, offset, limit)
	return
}

var rDAO IResourceDAO = &resourceDAO{}

func GetResourceDAO() IResourceDAO {
	return rDAO
}

func SetResourceDAO(v IResourceDAO) {
	rDAO = v
}
