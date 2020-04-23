package department

import (
	"DBproject1/controller"

	"github.com/labstack/echo"
)

type saveDepartmentReq struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

type updateDepartmentReq struct {
	ID          uint64 `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

type viewRequest struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

type selectDepartmentReq struct {
	ID uint64 `json:"ID"`
}

func SaveDepartment(e echo.Context) error {
	return controller.ExecHandler(e, &saveDepartmentReq{}, saveDepartment)
}

func LoadUpdateDepartmentForm(e echo.Context) error {
	return controller.ExecHandler(e, &selectDepartmentReq{}, loadUpdateDepartmentForm)
}

func UpdateDepartment(e echo.Context) error {
	return controller.ExecHandler(e, &updateDepartmentReq{}, updateDepartment)
}

func ViewDepartmentMember(e echo.Context) error {
	return controller.ExecHandler(e, &selectDepartmentReq{}, viewDepartmentMember)
}

func ViewDepartmentPlan(e echo.Context) error {
	return controller.ExecHandler(e, &selectDepartmentReq{}, viewDepartmentPlan)
}

func SelectAllDepartment(e echo.Context) error {
	return controller.ExecHandler(e, nil, selectAllDepartment)
}

func SelectDepartmentByView(e echo.Context) error {
	return controller.ExecHandler(e, &viewRequest{}, selectDepartmentByView)
}
