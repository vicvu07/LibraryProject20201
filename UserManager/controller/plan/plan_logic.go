package plan

import (
	"DBproject1/core"
	"DBproject1/dao"
	"DBproject1/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func savePlan(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*savePlanReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Save Plan", Data: req.Name}

	db := core.GetDB()
	pDAO := dao.GetPlanDAO()

	fatherPlanID, _ := strconv.ParseUint(req.FatherPlanID, 10, 64)

	er := pDAO.SavePlan(ctx, db, req.Name, req.Description, fatherPlanID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, nil, lg, false, nil
}

func donePlan(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*selectPlanReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Done Plan", Data: req.ID}

	db := core.GetDB()
	pDAO := dao.GetPlanDAO()

	er := pDAO.DonePlan(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, nil, lg, false, nil
}

func updatePlan(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*updatePlanReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Update Plan Detail", Data: req.ID}

	db := core.GetDB()
	pDAO := dao.GetPlanDAO()

	fatherPlanID, _ := strconv.ParseUint(req.FatherPlanID, 10, 64)

	er := pDAO.UpdatePlan(ctx, db, req.ID, req.Name, req.Description, fatherPlanID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, nil, lg, false, nil
}

func updatePlanUser(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*updatePlanUserReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Update Plan User", Data: req.ID}

	db := core.GetDB()
	pDAO := dao.GetPlanDAO()

	fmt.Println(req.User)
	fmt.Println(req.ID)

	er := pDAO.DeleteAllFromPlan(ctx, db, req.ID, "I")
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	for i := 0; i < len(req.User); i++ {
		er := pDAO.AddToPlan(ctx, db, "I", req.User[i].ID, req.ID)
		if er != nil {
			statusCode, err = http.StatusInternalServerError, er
			return
		}
	}

	return http.StatusOK, nil, lg, false, nil
}

func updatePlanDepartment(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*updatePlanDepartmentReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Update Plan Department", Data: req.ID}

	db := core.GetDB()
	pDAO := dao.GetPlanDAO()

	er := pDAO.DeleteAllFromPlan(ctx, db, req.ID,"D")
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	for i := 0; i < len(req.Department); i++ {
		er := pDAO.AddToPlan(ctx, db, "D", req.Department[i].ID, req.ID)
		if er != nil {
			statusCode, err = http.StatusInternalServerError, er
			return
		}
	}

	return http.StatusOK, nil, lg, false, nil
}

func selectAllFinishedPlan(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	ctx := c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Select All Finished Plan", Data: ""}

	db := core.GetDB()
	pDAO := dao.GetPlanDAO()

	result, er := pDAO.SelectAllFinishedPlan(ctx, db)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, result, lg, false, nil
}

func selectAllOnGoingPlan(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	ctx := c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Select All On Going Plan", Data: ""}

	db := core.GetDB()
	pDAO := dao.GetPlanDAO()

	result, er := pDAO.SelectAllOnGoingPlan(ctx, db)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, result, lg, false, nil
}

func selectFinishedPlan(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*selectPlanReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Show Finished Plan", Data: req.ID}

	db := core.GetDB()
	pDAO := dao.GetPlanDAO()

	result, er := pDAO.SelectChildrenDonePlan(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, result, lg, false, nil
}

func selectOnGoingPlan(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*selectPlanReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Show On Going Plan", Data: req.ID}

	db := core.GetDB()
	pDAO := dao.GetPlanDAO()

	result, er := pDAO.SelectChildrenOnGoingPlan(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, result, lg, false, nil
}

func selectPlanByID(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*selectPlanReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Show Plan Detail", Data: req.ID}

	db := core.GetDB()
	pDAO := dao.GetPlanDAO()

	result, er := pDAO.SelectPlanByID(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	var plan = &model.PlanForUI{
		ID:           result.ID,
		Name:         result.Name,
		Description:  result.Description,
		FatherPlanID: strconv.Itoa(result.FatherPlanID),
	}

	return http.StatusOK, plan, lg, false, nil
}

func selectPlanUser(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*selectPlanReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Show Plan User", Data: req.ID}

	db := core.GetDB()
	pDAO := dao.GetPlanDAO()

	result, er := pDAO.SelectUserInPlan(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, result, lg, false, nil
}

func selectPlanDepartment(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*selectPlanReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Show Plan Department", Data: req.ID}

	db := core.GetDB()
	pDAO := dao.GetPlanDAO()

	result, er := pDAO.SelectDepartmentInPlan(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, result, lg, false, nil
}
