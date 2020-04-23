package dao

import (
	"DBproject1/core"
	"DBproject1/model"
	"context"

	"github.com/linxGnu/mssqlx"
)

const (
	sqlActionSave         = "INSERT INTO action(name,description) VALUES (?,?)"
	sqlActionDelete       = "DELETE FROM action WHERE name=?"
	sqlActionSelectAll    = "SELECT *FROM action"
	sqlActionSelectByView = "SELECT *FROM action ORDER BY id DESC LIMIT ?, ?"
)

type IActionDAO interface {
	SaveAction(ctx context.Context, db *mssqlx.DBs, name, desc string) (err error)
	DeleteAction(ctx context.Context, db *mssqlx.DBs, name string) (err error)
	SelectAllAction(ctx context.Context, db *mssqlx.DBs) (result []*model.Action, err error)
	SelectActionByView(ctx context.Context, db *mssqlx.DBs, offset, limit int64) (result []*model.Action, err error)
}

type actionDAO struct{}

func (c *actionDAO) SaveAction(ctx context.Context, db *mssqlx.DBs, name, desc string) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}
	_, err = db.ExecContext(ctx, sqlActionSave, name, desc)
	return
}

func (c *actionDAO) DeleteAction(ctx context.Context, db *mssqlx.DBs, name string) (err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}
	_, err = db.ExecContext(ctx, sqlActionDelete, name)
	return
}

func (c *actionDAO) SelectAllAction(ctx context.Context, db *mssqlx.DBs) (result []*model.Action, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}
	err = db.SelectContext(ctx, &result, sqlActionSelectAll)
	return
}

func (c *actionDAO) SelectActionByView(ctx context.Context, db *mssqlx.DBs, offset, limit int64) (result []*model.Action, err error) {
	// Validate input
	if db == nil {
		err = core.ErrDBObjNull
		return
	}

	err = db.SelectContext(ctx, &result, sqlActionSelectByView, offset, limit)
	return
}

var aDAO IActionDAO = &actionDAO{}

func GetActionDAO() IActionDAO {
	return aDAO
}

func SetActionDAO(v IActionDAO) {
	aDAO = v
}
