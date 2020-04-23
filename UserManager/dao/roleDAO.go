package dao

import (
	"DBproject1/core"
	"DBproject1/model"
	"context"

	"github.com/linxGnu/mssqlx"
)

const (
	sqlRoleSave         = "INSERT INTO role(name,description) VALUES (?,?)"
	sqlRoleDelete       = "DELETE FROM role WHERE name=?"
	sqlRoleSelectAll    = "SELECT *FROM role"
	sqlRoleSelectByView = "SELECT *FROM role ORDER BY id DESC LIMIT ?, ?"
)

type IRoleDAO interface {
	SaveRole(ctx context.Context, db *mssqlx.DBs, name, desc string) (err error)
	DeleteRole(ctx context.Context, db *mssqlx.DBs, name string) (err error)
	SelectAllRole(ctx context.Context, db *mssqlx.DBs) (result []*model.Role, err error)
	SelectRoleByView(ctx context.Context, db *mssqlx.DBs, offset, limit int64) (result []*model.Role, err error)
}

type roleDAO struct{}

func (c *roleDAO) SaveRole(ctx context.Context, db *mssqlx.DBs, name string, desc string) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}
	_, err = db.ExecContext(ctx, sqlRoleSave, name, desc)
	return
}

func (c *roleDAO) DeleteRole(ctx context.Context, db *mssqlx.DBs, name string) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}
	_, err = db.ExecContext(ctx, sqlRoleDelete, name)
	return
}

func (c *roleDAO) SelectAllRole(ctx context.Context, db *mssqlx.DBs) (result []*model.Role, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}
	err = db.SelectContext(ctx, &result, sqlRoleSelectAll)
	return
}

func (c *roleDAO) SelectRoleByView(ctx context.Context, db *mssqlx.DBs, offset, limit int64) (result []*model.Role, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlRoleSelectByView, offset, limit)
	return
}

var roDAO IRoleDAO = &roleDAO{}

func GetRoleDAO() IRoleDAO {
	return roDAO
}

func SetRoleDAO(v IRoleDAO) {
	roDAO = v
}
