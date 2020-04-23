package action

import (
	"DBproject1/core"
	"DBproject1/dao"
	"DBproject1/model"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func save(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	// Log login info
	req := request.(*action)
	ctx := c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Add resource", Data: req.Name}
	db := core.GetDB()
	aDAO := dao.GetActionDAO()
	er := aDAO.SaveAction(ctx, db, req.Name, req.Description)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	return http.StatusOK, nil, lg, false, nil
}

func delete(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	// Log login info
	req := request.(*action)
	ctx := c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Delete resource", Data: req.Name}
	db := core.GetDB()
	aDAO := dao.GetActionDAO()
	e := core.GetCasbinEnforcer()

	er := aDAO.DeleteAction(ctx, db, req.Name)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	perm := e.GetPolicy()
	for i := 0; i < len(perm); i++ {
		if perm[i][2] == req.Name {
			erro := e.RemovePolicy(perm[i][0], perm[i][1], perm[i][2])
			if !erro {
				statusCode, err = http.StatusInternalServerError, fmt.Errorf("Can't delete action")
				return
			}
		}
	}

	return http.StatusOK, nil, lg, false, nil
}

func selectByView(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	// Log login info
	req := request.(*viewRequest)
	ctx := c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Select Acion By View", Data: req.Limit}
	db := core.GetDB()
	aDAO := dao.GetActionDAO()
	res, er := aDAO.SelectActionByView(ctx, db, req.Offset, req.Limit)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	return http.StatusOK, res, lg, false, nil
}

func selectAll(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	ctx := c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Select All Action", Data: "selectAll"}
	db := core.GetDB()
	aDAO := dao.GetActionDAO()
	res, er := aDAO.SelectAllAction(ctx, db)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	var resp []string
	for i := 0; i < len(res); i++ {
		resp = append(resp, res[i].Name)
	}
	fmt.Println(resp)
	return http.StatusOK, resp, lg, false, nil
}
