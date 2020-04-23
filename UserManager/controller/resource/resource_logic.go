package resource

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
	req := request.(*resource)
	ctx := c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Add resource", Data: req.Name}
	db := core.GetDB()
	rDAO := dao.GetResourceDAO()
	er := rDAO.SaveResource(ctx, db, req.Name, req.Description)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	return http.StatusOK, nil, lg, false, nil
}

func delete(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	// Log login info
	req := request.(*resource)
	ctx := c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Delete resource", Data: req.Name}
	db := core.GetDB()
	rDAO := dao.GetResourceDAO()
	e := core.GetCasbinEnforcer()
	er := rDAO.DeleteResource(ctx, db, req.Name)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	perm := e.GetPolicy()
	for i := 0; i < len(perm); i++ {
		if perm[i][1] == req.Name {
			erro := e.RemovePolicy(perm[i][0], perm[i][1], perm[i][2])
			if !erro {
				statusCode, err = http.StatusInternalServerError, fmt.Errorf("Can't delete resource")
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
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Select By Resource View", Data: req.Limit}
	db := core.GetDB()
	rDAO := dao.GetResourceDAO()
	res, er := rDAO.SelectResourceByView(ctx, db, req.Offset, req.Limit)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	fmt.Println(res)
	return http.StatusOK, res, lg, false, nil
}

func selectAll(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	ctx := c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Select All Resource", Data: "selectAll"}
	db := core.GetDB()
	rDAO := dao.GetResourceDAO()
	res, er := rDAO.SelectAllResource(ctx, db)
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
