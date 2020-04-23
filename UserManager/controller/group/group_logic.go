package group

import (
	"DBproject1/core"
	"DBproject1/dao"
	"DBproject1/model"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func createGroup(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req := request.(*reqGroup)
	ctx := c.Request().Context()
	// Log login info
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Create Group", Data: req.GroupName}
	db := core.GetDB()

	roleDAO := dao.GetRoleDAO()
	er := roleDAO.SaveRole(ctx, db, req.GroupName, req.Description)
	if er != nil {
		return http.StatusInternalServerError, nil, lg, false, er
	}

	return http.StatusOK, nil, lg, false, nil
}

func deleteGroup(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req := request.(*reqGroup)
	ctx := c.Request().Context()
	// Log login info
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Delete Group", Data: req.GroupName}
	db, e := core.GetDB(), core.GetCasbinEnforcer()
	roleDAO := dao.GetRoleDAO()
	//delete from db
	er := roleDAO.DeleteRole(ctx, db, req.GroupName)
	if er != nil {
		return http.StatusInternalServerError, nil, lg, false, er
	}
	// delete role
	e.DeleteRole(req.GroupName)
	// delete policies added to group
	perm := e.GetPermissionsForUser(req.GroupName)
	if len(perm) > 0 {
		for i := 0; i < len(perm); i++ {
			erro := e.RemovePolicy(perm[i][0], perm[i][1], perm[i][2])
			if !erro {
				statusCode, err = http.StatusInternalServerError, fmt.Errorf("Can't remove group permission")
				return
			}
		}
	}

	// remove user out of group
	user := e.GetUsersForRole(req.GroupName)
	if len(user) > 0 {
		for i := 0; i < len(user); i++ {
			erro := e.DeleteRoleForUser(user[i], req.GroupName)
			if !erro {
				statusCode, err = http.StatusInternalServerError, fmt.Errorf("Can't remove user out of group %s", req.GroupName)
				return
			}
		}
	}

	return http.StatusOK, nil, lg, false, nil
}

func loadGroupPermission(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req := request.(*reqLoadGroupPermission)
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Load Update Permission", Data: req.GroupName}
	e := core.GetCasbinEnforcer()
	//	var tresource, taction []string

	r := e.GetPermissionsForUser(req.GroupName)
	// for i := 0; i < len(r); i++ {
	// 	tresource = append(tresource, r[i][1])
	// 	taction = append(taction, r[i][2])
	// }
	fmt.Println(req.GroupName)
	log.Println(r)

	var tPerm [][]string
	for i := 0; i < len(r); i++ {
		if r[i][0] == req.GroupName {
			var tmp []string
			tmp = append(tmp, r[i][1])
			tmp = append(tmp, r[i][2])
			tPerm = append(tPerm, tmp)
		}
	}
	log.Println(tPerm)
	res := &model.PermissionForUI{
		TPerm: tPerm,
	}

	return http.StatusOK, res, lg, false, nil
}

func updateGroupPermission(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req := request.(*reqGroupPermission)
	// Log login info
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Update Group Permission", Data: req.GroupName}
	e := core.GetCasbinEnforcer()

	per := e.GetPermissionsForUser(req.GroupName)
	log.Println(per)
	if len(per) > 0 {
		erro := e.DeletePermissionsForUser(req.GroupName)
		if !erro {
			statusCode, err = http.StatusInternalServerError, fmt.Errorf("Horus: Can't delete user permissions")
			return
		}
	}
	// create new
	for i := 0; i < len(req.Permission); i++ {
		fmt.Println(req.Permission[i][0], req.Permission[i][1])
		erro := e.AddPolicy(req.GroupName, req.Permission[i][0], req.Permission[i][1])
		if !erro {
			statusCode, err = http.StatusInternalServerError, fmt.Errorf("Horus: Can't add group permissions")
			return
		}
	}
	return http.StatusOK, nil, lg, false, nil
}

func selectGroupByView(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req := request.(*viewRequest)
	// Log login info
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "selectGroupByView", Data: req.Limit}
	db := core.GetDB()
	ctx := c.Request().Context()

	roDAO := dao.GetRoleDAO()
	res, er := roDAO.SelectRoleByView(ctx, db, req.Offset, req.Limit)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	return http.StatusOK, res, lg, false, nil
}

func selectAllGroup(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Select All Groups", Data: "select all"}
	db := core.GetDB()
	ctx := c.Request().Context()

	roDAO := dao.GetRoleDAO()
	res, er := roDAO.SelectAllRole(ctx, db)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	var resp []string
	for i := 0; i < len(res); i++ {
		resp = append(resp, res[i].Name)
	}
	log.Println(resp)
	return http.StatusOK, resp, lg, false, nil
}
