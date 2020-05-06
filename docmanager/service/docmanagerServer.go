package service

import (
	"context"

	"github.com/pinezapple/LibraryProject20201/docmanager/core"
	"github.com/pinezapple/LibraryProject20201/docmanager/dao"
	"github.com/pinezapple/LibraryProject20201/skeleton/logger"
	"github.com/pinezapple/LibraryProject20201/skeleton/model"
	"github.com/pinezapple/LibraryProject20201/skeleton/model/docmanagerModel"
)

type docmanagerServer struct {
	lg *model.LogFormat
}

func (d *docmanagerServer) SelectAllDoc(ctx context.Context, req *docmanagerModel.SelectAllDocReq) (resp *docmanagerModel.SelectAllDocResp, err error) {
	logger.LogInfo(d.lg, "rpc Select all Doc Req")

	docs, err := dao.GetDocDAO().SelectAllDoc(ctx, core.GetDB())
	if err != nil {
		logger.LogErr(d.lg, err)
		return &docmanagerModel.SelectAllDocResp{Code: 1, Message: err.Error()}, nil
	}

	logger.LogInfo(d.lg, "rpc Select all Doc Success")
	return &docmanagerModel.SelectAllDocResp{Code: 0, Documents: docs}, nil
}

func (d *docmanagerServer) SelectDocByID(ctx context.Context, req *docmanagerModel.SelectDocByIDReq) (resp *docmanagerModel.SelectDocByIDResp, err error) {
	logger.LogInfo(d.lg, "rpc Select Doc By ID Req")

	doc, err := dao.GetDocDAO().SelectDocByID(ctx, core.GetDB(), req.DocID)
	if err != nil {
		logger.LogErr(d.lg, err)
		return &docmanagerModel.SelectDocByIDResp{Code: 1, Message: err.Error()}, nil
	}
	logger.LogInfo(d.lg, "rpc Select Doc By ID Success")
	return &docmanagerModel.SelectDocByIDResp{Code: 0, Documents: doc}, nil
}

func (d *docmanagerServer) SaveDoc(ctx context.Context, req *docmanagerModel.SaveDocReq) (resp *docmanagerModel.SaveDocResp, err error) {
	logger.LogInfo(d.lg, "rpc Save Doc Req")

	err = dao.GetDocDAO().SaveDoc(ctx, core.GetDB(), req.Doc)
	if err != nil {
		logger.LogErr(d.lg, err)
		return &docmanagerModel.SaveDocResp{Code: 1, Message: err.Error()}, nil
	}
	logger.LogInfo(d.lg, "rpc Select Doc By ID Success")
	return &docmanagerModel.SaveDocResp{Code: 0}, nil
}

func (d *docmanagerServer) UpdateDoc(ctx context.Context, req *docmanagerModel.UpdateDocReq) (resp *docmanagerModel.UpdateDocResp, err error) {
	logger.LogInfo(d.lg, "rpc Update Doc Req")

	err = dao.GetDocDAO().UpdateDoc(ctx, core.GetDB(), req.Doc)
	if err != nil {
		logger.LogErr(d.lg, err)
		return &docmanagerModel.UpdateDocResp{Code: 1, Message: err.Error()}, nil
	}
	logger.LogInfo(d.lg, "rpc Select Doc By ID Success")
	return &docmanagerModel.UpdateDocResp{Code: 0}, nil
}

func (d *docmanagerServer) DeleteDoc(ctx context.Context, req *docmanagerModel.DeleteDocReq) (resp *docmanagerModel.DeleteDocResp, err error) {
	logger.LogInfo(d.lg, "rpc Delete Doc Req")

	err = dao.GetDocDAO().DelDoc(ctx, core.GetDB(), req.DocID)
	if err != nil {
		logger.LogErr(d.lg, err)
		return &docmanagerModel.DeleteDocResp{Code: 1, Message: err.Error()}, nil
	}
	logger.LogInfo(d.lg, "rpc Select Doc By ID Success")
	return &docmanagerModel.DeleteDocResp{Code: 0}, nil
}

// -------------------------------------------- Borrow Form ------------------------------------------------------
func (d *docmanagerServer) SaveBorrowForm(ctx context.Context, req *docmanagerModel.SaveBorrowFormReq) (resp *docmanagerModel.SaveBorrowFormResp, err error) {
	logger.LogInfo(d.lg, "rpc Save Borrow Form Req")

	err = dao.GetDocDAO().SaveBorrowForm(ctx, core.GetDB(), req.Borrowform)
	if err != nil {
		logger.LogErr(d.lg, err)
		return &docmanagerModel.SaveBorrowFormResp{Code: 1, Message: err.Error()}, nil
	}
	logger.LogInfo(d.lg, "rpc Save Borrow Form Success")
	return &docmanagerModel.SaveBorrowFormResp{Code: 0}, nil
}

func (d *docmanagerServer) UpdateBorrowFormStatus(ctx context.Context, req *docmanagerModel.UpdateBorrowFormStatusReq) (resp *docmanagerModel.UpdateBorrowFormStatusResp, err error) {
	logger.LogInfo(d.lg, "rpc Update Borrow Form Req")

	err = dao.GetDocDAO().UpdateBorrowFormStatus(ctx, core.GetDB(), req.FormID, int(req.Status))
	if err != nil {
		logger.LogErr(d.lg, err)
		return &docmanagerModel.UpdateBorrowFormStatusResp{Code: 1, Message: err.Error()}, nil
	}
	logger.LogInfo(d.lg, "rpc Update Borrow Form Success")
	return &docmanagerModel.UpdateBorrowFormStatusResp{Code: 0}, nil

}

func (d *docmanagerServer) SelectBorrowFormByID(ctx context.Context, req *docmanagerModel.SelectBorrowFormByIDReq) (resp *docmanagerModel.SelectBorrowFormByIDResp, err error) {
	logger.LogInfo(d.lg, "rpc Select Borrow Form Req")

	res, err := dao.GetDocDAO().SelectBorrowFormByID(ctx, core.GetDB(), req.FormID)
	if err != nil {
		logger.LogErr(d.lg, err)
		return &docmanagerModel.SelectBorrowFormByIDResp{Code: 1, Message: err.Error()}, nil
	}
	logger.LogInfo(d.lg, "rpc Select Borrow Form Success")
	return &docmanagerModel.SelectBorrowFormByIDResp{Code: 0, Borrowform: res}, nil
}
