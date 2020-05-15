package dao

import (
	"context"
	"time"

	"github.com/linxGnu/mssqlx"
	"github.com/pinezapple/LibraryProject20201/docmanager/core"
	"github.com/pinezapple/LibraryProject20201/skeleton/model/docmanagerModel"
)

const (
	sqlSelectAll              = "SELECT * FROM doc"
	sqlSelectByID             = "SELECT * FROM doc WHERE id_doc = ?"
	sqlSave                   = "INSERT INTO doc(id_doc,doc_name,doc_author,doc_type,doc_description,status,fee) VALUES (?,?,?,?,?,?,?)"
	sqlDelete                 = "DELETE FROM doc WHERE id_doc = ?"
	sqlUpdateStatus           = "UPDATE doc SET status = ?, id_borrow = ?,updated_at = ? WHERE id_doc= ?"
	sqlUpdate                 = "UPDATE doc SET doc_name = ?, doc_author = ?, doc_type =?, doc_description = ?, status = ?, fee = ?, updated_at = ? WHERE id_doc = ?"
	sqlSaveBorrowForm         = "INSERT INTO borrowform(id_borrow, id_doc, id_cus, id_lib, status, start_at, end_at) VALUE (?,?,?,?,?,?,?)"
	sqlUpdateBorrowFormStatus = "UPDATE borrowform SET status = ?, updated_at = ? WHERE id_borrow = ?"
	sqlSelectBorrowFormByID   = "SELECT * FROM borrowform WHERE id_borrow = ?"
	sqlSelecetIdDoc           = "SELECT id_doc FROM doc WHERE id_borrow = ?"
)

type IDocDAO interface {
	SelectAllDoc(ctx context.Context, db *mssqlx.DBs) (result []*docmanagerModel.Doc, err error)
	SelectDocByID(ctx context.Context, db *mssqlx.DBs, id uint64) (result *docmanagerModel.Doc, err error)
	SaveDoc(ctx context.Context, db *mssqlx.DBs, doc *docmanagerModel.Doc) (err error)
	UpdateDoc(ctx context.Context, db *mssqlx.DBs, doc *docmanagerModel.Doc) (err error)
	DelDoc(ctx context.Context, db *mssqlx.DBs, id uint64) (err error)
	//--------------- BorrowForm --------------
	SaveBorrowForm(ctx context.Context, db *mssqlx.DBs, form *docmanagerModel.BorrowForm) (err error)
	UpdateBorrowFormStatus(ctx context.Context, db *mssqlx.DBs, id uint64, status int) (err error)
	SelectBorrowFormByID(ctx context.Context, db *mssqlx.DBs, id uint64) (result *docmanagerModel.BorrowForm, err error)
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

	_, err = db.ExecContext(ctx, sqlSave, doc.ID, doc.Name, doc.Author, doc.Type, doc.Descriptor, doc.Status, doc.Fee)
	return
}

func (d *docDAO) UpdateDoc(ctx context.Context, db *mssqlx.DBs, doc *docmanagerModel.Doc) (err error) {
	if db == nil {
		return core.ErrDBObjNull
	}

	_, err = db.ExecContext(ctx, sqlUpdate, doc.Name, doc.Author, doc.Type, doc.Descriptor, doc.Status, doc.Fee, time.Now(), doc.ID)
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

//-------------------------- Borrow Form ---------------------------
func (d *docDAO) SaveBorrowForm(ctx context.Context, db *mssqlx.DBs, form *docmanagerModel.BorrowForm) (err error) {
	if db == nil {
		return core.ErrDBObjNull
	}

	_, err = db.ExecContext(ctx, sqlSaveBorrowForm, form.ID, form.DocID, form.CusID, form.LibID, form.Status, form.StartAt, form.EndAt)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sqlUpdateStatus, form.Status, form.ID, time.Now(), form.DocID)
	return
}

func (d *docDAO) UpdateBorrowFormStatus(ctx context.Context, db *mssqlx.DBs, id uint64, status int) (err error) {
	if db == nil {
		return core.ErrDBObjNull
	}
	_, err = db.ExecContext(ctx, sqlUpdateBorrowFormStatus, status, time.Now(), id)
	if err != nil {
		return err
	}
	var id_doc uint64
	err = db.SelectContext(ctx, &id_doc, sqlSelecetIdDoc, id)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, sqlUpdateStatus, status, id, time.Now(), id_doc)
	return
}

func (d *docDAO) SelectBorrowFormByID(ctx context.Context, db *mssqlx.DBs, id uint64) (result *docmanagerModel.BorrowForm, err error) {
	if db == nil {
		return nil, core.ErrDBObjNull
	}

	err = db.SelectContext(ctx, result, sqlSelectBorrowFormByID, id)
	return
}

var dDao IDocDAO = &docDAO{}

func SetDocDAO(d IDocDAO) {
	dDao = d
}

func GetDocDAO() (d IDocDAO) {
	return dDao
}
