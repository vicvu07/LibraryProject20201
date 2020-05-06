package service

import (
	"context"

	"github.com/pinezapple/LibraryProject20201/skeleton/model"
	"github.com/pinezapple/LibraryProject20201/skeleton/model/docmanagerModel"
)

type docmanagerServer struct {
	lg *model.LogFormat
}

func (d *docmanagerServer) SelectAllDoc(ctx context.Context, req *docmanagerModel.SelectAllDocReq) (resp *docmanagerModel.SelectAllDocResp, err error) {
	return
}

func (d *docmanagerServer) SelectDocByID(ctx context.Context, req *docmanagerModel.SelectDocByIDReq) (resp *docmanagerModel.SelectDocByIDResp, err error) {
	return
}

func (d *docmanagerServer) SaveDoc(ctx context.Context, req *docmanagerModel.SaveDocReq) (resp *docmanagerModel.SaveDocResp, err error) {
	return
}

func (d *docmanagerServer) UpdateDoc(ctx context.Context, req *docmanagerModel.UpdateDocReq) (resp *docmanagerModel.UpdateDocResp, err error) {
	return
}
func (d *docmanagerServer) DeleteDoc(ctx context.Context, req *docmanagerModel.DeleteDocReq) (resp *docmanagerModel.DeleteDocResp, err error) {
	return
}
