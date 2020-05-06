package dao

import (
	"context"

	"github.com/linxGnu/mssqlx"
	"github.com/pinezapple/LibraryProject20201/docmanager/core"
	"github.com/pinezapple/LibraryProject20201/skeleton/model/docmanagerModel"
)

const (
	sqlSelectAll    = "SELECT * FROM doc"
	sqlSelectByID   = "SELECT * FROM doc WHERE doc_id = ?"
	sqlSave         = "INSERT INTO doc(doc_name,doc_author,doc_type,doc_description,status,fee) VALUES (?,?,?,?,?,?)"
	sqlDelete       = "DELETE FROM doc WHERE doc_id = ?"
	sqlUpdateStatus = "UPDATE doc SET status = ?, updated_at = ? WHERE doc_id= ?"
	sqlUpdate       = "UPDATE doc SET doc_name = ?, doc_author = ?, doc_type =?, doc_description = ?, status = ?, fee = ? WHERE doc_id = ?"
)

type IDocDAO interface {
	SelectAllDoc(ctx context.Context, db *mssqlx.DBs) (result []*docmanagerModel.Doc, err error)
	SelectDocByID(ctx context.Context, db *mssqlx.DBs, id uint64) (result *docmanagerModel.Doc, err error)
	SaveDoc(ctx context.Context, db *mssqlx.DBs, doc *docmanagerModel.Doc) (err error)
	UpdateDoc(ctx context.Context, db *mssqlx.DBs, doc *docmanagerModel.Doc) (err error)
	DelDoc(ctx context.Context, db *mssqlx.DBs, id uint64) (err error)
}

type docDAO struct {
}

func (d *docDAO) SelectAllDoc(ctx context.Context, db *mssqlx.DBs) (result []*docmanagerModel.Doc, err error) {
	if db == nil {
		return nil, core.ErrDBObjNull
	}

	err = db.SelectContext(ctx, &result, sqlSelectAll)
	return
}

func (d *docDAO) SelectDocByID(ctx context.Context, db *mssqlx.DBs, id uint64) (result *docmanagerModel.Doc, err error) {
	if db == nil {
		return nil, core.ErrDBObjNull
	}

	err = db.SelectContext(ctx, result, sqlSelectByID, id)
	return
}

func (d *docDAO) SaveDoc(ctx context.Context, db *mssqlx.DBs, doc *docmanagerModel.Doc) (err error) {
	if db == nil {
		return core.ErrDBObjNull
	}

	_, err = db.ExecContext(ctx, sqlSave, doc.Name, doc.Author, doc.Type, doc.Descriptor, doc.Status, doc.Fee)
	return
}

func (d *docDAO) UpdateDoc(ctx context.Context, db *mssqlx.DBs, doc *docmanagerModel.Doc) (err error) {
	if db == nil {
		return core.ErrDBObjNull
	}

	_, err = db.ExecContext(ctx, sqlUpdate, doc.Name, doc.Author, doc.Type, doc.Descriptor, doc.Status, doc.Fee, doc.ID)
	return
}

func (d *docDAO) DelDoc(ctx context.Context, db *mssqlx.DBs, id uint64) (err error) {
	if db == nil {
		return core.ErrDBObjNull
	}

	_, err = db.ExecContext(ctx, sqlDelete, id)
	return
}

func (d *docDAO) UpdateStatus(ctx context.Context, db *mssqlx.DBs, id uint64, status int) (err error) {
	if db == nil {
		return core.ErrDBObjNull
	}

	_, err = db.ExecContext(ctx, sqlUpdateStatus, status, id)
	return
}

var dDao IDocDAO = &docDAO{}

func SetDocDAO(d IDocDAO) {
	dDao = d
}

func GetDocDAO() (d IDocDAO) {
	return dDao
}
