package plan

import (
	"DBproject1/controller"
	"DBproject1/model"

	"github.com/labstack/echo"
)

type savePlanReq struct {
	Name         string `json:"Name"`
	Description  string `json:"Description"`
	FatherPlanID string `json:"FatherPlanID"`
}

type updatePlanReq struct {
	ID           uint64 `json:"ID"`
	Name         string `json:"Name"`
	Description  string `json:"Description"`
	FatherPlanID string `json:"FatherPlanID"`
}

type viewRequest struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

type selectPlanReq struct {
	ID uint64 `json:"ID"`
}

type updatePlanUserReq struct {
	ID   uint64                    `json:"ID"`
	User []model.UserDetailForPlan `json:"User"`
}

type updatePlanDepartmentReq struct {
	ID         uint64                    `json:"ID"`
	Department []model.DepartmentForPlan `json:"Department"`
}

func SavePlan(e echo.Context) error {
	return controller.ExecHandler(e, &savePlanReq{}, savePlan)
}

func DonePlan(e echo.Context) error {
	return controller.ExecHandler(e, &selectPlanReq{}, donePlan)
}

func UpdatePlan(e echo.Context) error {
	return controller.ExecHandler(e, &updatePlanReq{}, updatePlan)
}

func UpdatePlanUser(e echo.Context) error {
	return controller.ExecHandler(e, &updatePlanUserReq{}, updatePlanUser)
}

func UpdatePlanDepartment(e echo.Context) error {
	return controller.ExecHandler(e, &updatePlanDepartmentReq{}, updatePlanDepartment)
}

func SelectAllFinishedPlan(e echo.Context) error {
	return controller.ExecHandler(e, nil, selectAllFinishedPlan)
}

func SelectAllOnGoingPlan(e echo.Context) error {
	return controller.ExecHandler(e, nil, selectAllOnGoingPlan)
}

func SelectFinishedPlan(e echo.Context) error {
	return controller.ExecHandler(e, &selectPlanReq{}, selectFinishedPlan)
}

func SelectOnGoingPlan(e echo.Context) error {
	return controller.ExecHandler(e, &selectPlanReq{}, selectOnGoingPlan)
}

func SelectPlanByID(e echo.Context) error {
	return controller.ExecHandler(e, &selectPlanReq{}, selectPlanByID)
}

func SelectPlanUser(e echo.Context) error {
	return controller.ExecHandler(e, &selectPlanReq{}, selectPlanUser)
}

func SelectPlanDepartment(e echo.Context) error {
	return controller.ExecHandler(e, &selectPlanReq{}, selectPlanDepartment)
}
