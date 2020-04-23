package department

import (
	"DBproject1/core"
	"DBproject1/dao"
	"DBproject1/model"
	"net/http"
	"fmt"
	"github.com/labstack/echo"
)

func saveDepartment(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*saveDepartmentReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Save Department", Data: req.Name}

	db := core.GetDB()
	dDAO := dao.GetDepartmentDAO()
	fmt.Println(req.Name)

	er := dDAO.CreateNewDepartment(ctx, db, req.Name, req.Description, 0)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, nil, lg, false, nil
}

func loadUpdateDepartmentForm(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*selectDepartmentReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "load Deparment form", Data: req.ID}

	db := core.GetDB()
	dDAO := dao.GetDepartmentDAO()

	result, er := dDAO.SelectDepartmentByID(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, result, lg, false, nil
}

func updateDepartment(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*updateDepartmentReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "update Deparment ", Data: req.ID}

	db := core.GetDB()
	dDAO := dao.GetDepartmentDAO()

	er := dDAO.UpdateDepartment(ctx, db, req.ID, req.Name, req.Description, 0)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, nil, lg, false, nil
}

func selectDepartmentByView(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*viewRequest), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "view Deparment ", Data: req.Offset}

	db := core.GetDB()
	dDAO := dao.GetDepartmentDAO()

	result, er := dDAO.SelectDepartmentByView(ctx, db, req.Offset, req.Limit)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, result, lg, false, nil
}

func selectAllDepartment(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	ctx := c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "All Deparment ", Data: "view all"}

	db := core.GetDB()
	dDAO := dao.GetDepartmentDAO()

	result, er := dDAO.SelectAllDepartment(ctx, db)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, result, lg, false, nil
}

func viewDepartmentMember(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*selectDepartmentReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "view users ", Data: req.ID}

	db := core.GetDB()
	dDAO := dao.GetDepartmentDAO()
	
	fmt.Println("view Department Member")
	fmt.Println(req.ID)

	oldUser, er := dDAO.SelectAllUserUsedToWorkInDepartment(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	newUser, er := dDAO.SelectAllUserWorkingInDepartment(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	var result = &model.UserDepartment{
		NewUser: newUser,
		OldUser: oldUser,
	}

	return http.StatusOK, result, lg, false, nil
}

func viewDepartmentPlan(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*selectDepartmentReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "view plan", Data: req.ID}

	db := core.GetDB()
	dDAO := dao.GetDepartmentDAO()

	oldPlan, er := dDAO.SelectAllDonePlanInDepartment(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return 
	}

	newPlan, er := dDAO.SelectAllOnGoingPlanInDepartment(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	var result = &model.UserPlan{
        DonePlan: oldPlan,
        OnGoingPlan: newPlan,
	}

	return http.StatusOK, result, lg, false, nil
}

