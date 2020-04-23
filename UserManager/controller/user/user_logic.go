package user

import (
	"DBproject1/core"
	"DBproject1/dao"
	"DBproject1/model"
	"fmt"
	"log"
	"net/http"
	"strings"
	"strconv"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func deleteUser(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*UserDeleteReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Delete User", Data: req.Username}

	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	e, db, conf := core.GetCasbinEnforcer(), core.GetDB(), core.GetConfig()
	if db == nil || conf == nil {
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("Horus: Database connection / Config not initialized")
		return
	}

	userDAO := dao.GetUserDAO()

	er := userDAO.Delete(ctx, db, req.Username, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	per := e.GetPermissionsForUser(req.Username)
	log.Println(per)
	if len(per) > 0 {
		erro := e.DeletePermissionsForUser(req.Username)
		if !erro {
			statusCode, err = http.StatusInternalServerError, fmt.Errorf("Horus: Can't delete user permissions")
			return
		}
	}
	ro := e.GetRolesForUser(req.Username)
	if len(ro) > 0 {
		erro := e.DeleteRolesForUser(req.Username)
		if !erro {
			statusCode, err = http.StatusInternalServerError, fmt.Errorf("Horus: Can't move user out of groups")
			return
		}
	}

	return http.StatusOK, nil, lg, false, nil
}

func save(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	req, ctx := request.(*userSaveReq), c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Create User", Data: req.Username}

	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	db, conf := core.GetDB(), core.GetConfig()
	if db == nil || conf == nil {
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("Horus: Database connection / Config not initialized")
		return
	}

	salary, _ := strconv.ParseUint(req.Salary, 10, 64)
//	dp,_ := strconv.ParseUint(req.DepartmentID, 10, 64)
	userDAO := dao.GetUserDAO()
	_, er := userDAO.Create(ctx, db, conf.WebServer.Secure.SipHashSum0, conf.WebServer.Secure.SipHashSum1, req.Name, req.DOB, req.Sex, req.Position, req.PhoneNum, req.NationalID, salary, req.Username, req.DepartmentID, req.Password)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	
	dDAO := dao.GetDepartmentDAO()
	dp, er := dDAO.SelectDepartmentByID(ctx, db, req.DepartmentID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	
	er = dDAO.UpdateDepartment(ctx, db, req.DepartmentID, dp.Name, dp.Description, dp.TotalSalary + salary)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, nil, lg, false, nil
}

func selectUserByView(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	ctx := c.Request().Context()
	req := request.(*viewRequest)
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Select By View", Data: "Select By View"}

	db, conf := core.GetDB(), core.GetConfig()
	if db == nil || conf == nil {
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("Horus: Database connection / Config not initialized")
		return
	}

	userDAO := dao.GetUserDAO()
	users, er := userDAO.SelectView(ctx, db, req.Offset, req.Limit)
	//log.Println(users)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, users, lg, false, nil
}

func selectAllUser(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	ctx := c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Select All User", Data: "Select All"}

	db, conf := core.GetDB(), core.GetConfig()
	if db == nil || conf == nil {
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("Horus: Database connection / Config not initialized")
		return
	}

	userDAO := dao.GetUserDAO()
	users, er := userDAO.SelectAll(ctx, db)
	//log.Println(users)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, users, lg, false, nil
}

func changePassword(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logRespData bool, err error) {
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "horus.changePassword"}
	db := core.GetDB()
	req, ctx, userDAO, conf := request.(*reqChangePasswd), c.Request().Context(), dao.GetUserDAO(), core.GetConfig()

	user, er := userDAO.SelectByUsername(ctx, db, req.Username)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	} else if user == nil {
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("User not found")
		return
	}
	if !user.ValidateChecksum(conf.WebServer.Secure.SipHashSum0, conf.WebServer.Secure.SipHashSum1) {
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("Invalid user checksum")
		return
	}

	// TODO: using constants later
	if user.Status != 1 {
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("This account has been blocked. Please contact admin of this site")
		return
	}

	// Select user_sec by his username
	userSec, er := userDAO.SelectSecByUsername(ctx, db, req.Username)
	if err != nil {
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("User not found")
		return
	}
	if !userSec.ValidateChecksum(conf.WebServer.Secure.SipHashSum0, conf.WebServer.Secure.SipHashSum1) {
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("Invalid user sec checksum")
		return
	}

	if err = bcrypt.CompareHashAndPassword(userSec.Password, []byte(req.OldPassword)); err != nil {
		statusCode, err = http.StatusInternalServerError, fmt.Errorf("Old password incorrect")
		return
	}

	er = userDAO.UpdatePassword(ctx, db, conf.WebServer.Secure.SipHashSum0, conf.WebServer.Secure.SipHashSum1, req.Username, req.NewPassword)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, nil, lg, false, nil
}

func loadUpdateUserDetailForm(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logRespData bool, err error) {
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "horus.LoadUserUpdateUpdateForm"}
	req := request.(*LoadUserUpdateReq)
	ctx, userDAO := c.Request().Context(), dao.GetUserDAO()
	db := core.GetDB()

	user, er := userDAO.Select(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	var userUI = &model.UserDetailString{
		ID: user.ID,
		Name: user.Name,
		DOB: user.DOB,
		Sex: user.Sex,
		Position: user.Position,
		Salary: strconv.Itoa(user.Salary),
		PhoneNum: user.PhoneNum,
		NationalID: user.NationalID,
	}

	return http.StatusOK, userUI, lg, false, nil
}

func loadUpdateUserDeparmentForm(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logRespData bool, err error) {
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "horus.LoadUserDepartmentForm"}
	req := request.(*LoadUserUpdateReq)
	ctx, userDAO := c.Request().Context(), dao.GetUserDAO()
	db := core.GetDB()
	fmt.Println(req.ID)

	detail, er := userDAO.SelectUserDepartment(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	
	department, er := dao.GetDepartmentDAO().SelectDepartmentByID(ctx, db, detail.DepartmentID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	
	return http.StatusOK, department, lg, false, nil
}

// update User Detail
func updateUserDetail(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logRespData bool, err error) {
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "horus.UpdateUserDetail"}
	req := request.(*userUpdateReq)
	ctx, userDAO := c.Request().Context(), dao.GetUserDAO()
	db := core.GetDB()
	salary, _ := strconv.ParseUint(req.Salary, 10, 64)
	er := userDAO.UpdateUserDetail(ctx, db, req.ID, req.Name, req.DOB, req.Sex, req.Position, req.PhoneNum, req.NationalID, salary)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, nil, lg, false, nil
}

// update User Department
func updateUserDepartment(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logRespData bool, err error) {
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "horus.UpdateUserDepartment"}
	req := request.(*UpdateUserDepartmentReq)
	ctx, userDAO := c.Request().Context(), dao.GetUserDAO()
	db := core.GetDB()
	
	//id, _ := strconv.ParseUint(req.ID, 10, 64)
	//oldDepartmentID, _ := strconv.ParseUint(req.OldDepartmentID, 10, 64)
	//newDepartmentID, _ := strconv.ParseUint(req.NewDepartmentID, 10, 64)
//	fmt.Println(req.ID)
	//fmt.Println(oldDeparmentID)
//	fmt.Println(newDepartmentID)

	er := userDAO.UpdateUserDepartmentManagement(ctx, db, req.ID, req.OldDepartmentID, req.NewDepartmentID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	return http.StatusOK, nil, lg, false, nil
}

func loadUpdatePermissionForm(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logRespData bool, err error) {
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "horus.LoadUPermissionUpdateForm"}
	req := request.(*loadUpdateFormReq)
	e := core.GetCasbinEnforcer()
	//var taction, tresource []string

	// Get ticked resources and action
	r := e.GetPermissionsForUser(req.Username)
	// for i := 0; i < len(r); i++ {
	// 	tresource = append(tresource, r[i][1])
	// 	taction = append(taction, r[i][2])
	// }
	var tPerm [][]string
	for i := 0; i < len(r); i++ {
		var tmp []string
		tmp = append(tmp, r[i][1])
		tmp = append(tmp, r[i][2])
		tPerm = append(tPerm, tmp)
	}

	res := &model.PermissionForUI{
		TPerm: tPerm,
	}

	return http.StatusOK, res, lg, false, nil
}

func loadUpdateRoleForm(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logRespData bool, err error) {
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "horus.LoadRoleUpdateForm"}
	req := request.(*loadUpdateFormReq)
	e := core.GetCasbinEnforcer()
	// select all role

	tRole := e.GetRolesForUser(req.Username)
	res := &model.RoleForUI{
		TRole: tRole,
	}
	fmt.Println(res)
	return http.StatusOK, res, lg, false, nil
}

func updatePermission(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logRespData bool, err error) {
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "horus.updatePermission"}
	req := request.(*UpdatePermissionReq)

	e := core.GetCasbinEnforcer()
	// delete old
	per := e.GetPermissionsForUser(req.Username)
	log.Println(per)
	if len(per) > 0 {
		erro := e.DeletePermissionsForUser(req.Username)
		if !erro {
			statusCode, err = http.StatusInternalServerError, fmt.Errorf("Horus: Can't delete user permissions")
			return
		}
	}
	// create new
	for i := 0; i < len(req.Permission); i++ {
		erro := e.AddPolicy(req.Username, req.Permission[i][0], req.Permission[i][1])
		if !erro {
			statusCode, err = http.StatusInternalServerError, fmt.Errorf("Horus: Can't add user permissions")
			return
		}
	}

	return http.StatusOK, nil, lg, false, nil
}

func updateRole(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logRespData bool, err error) {
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "horus.updateRole"}
	req := request.(*UpdateRoleReq)

	e := core.GetCasbinEnforcer()
	// delete old
	ro := e.GetRolesForUser(req.Username)
	if len(ro) > 0 {
		erro := e.DeleteRolesForUser(req.Username)
		if !erro {
			statusCode, err = http.StatusInternalServerError, fmt.Errorf("Horus: Can't move user out of groups")
			return
		}
	}
	fmt.Println(req.Role)
	if len(req.Role) > 0 {
		// create new
		for i := 0; i < len(req.Role); i++ {
			erro := e.AddRoleForUser(req.Username, req.Role[i])
			if !erro {
				statusCode, err = http.StatusInternalServerError, fmt.Errorf("Horus: Can't add user to groups")
				return
			}
		}
	}

	return http.StatusOK, nil, lg, false, nil
}

func selectAllPermission(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	ctx := c.Request().Context()
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "Select All Permission", Data: "selectAll"}
	db := core.GetDB()
	aDAO := dao.GetActionDAO()
	a, er := aDAO.SelectAllAction(ctx, db)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	var action []string
	for i := 0; i < len(a); i++ {
		action = append(action, a[i].Name)
	}

	rDAO := dao.GetResourceDAO()
	r, er := rDAO.SelectAllResource(ctx, db)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	var resource []string
	for i := 0; i < len(r); i++ {
		resource = append(resource, r[i].Name)
	}

	var resp [][]string
	for i := 0; i < len(resource); i++ {
		for j := 0; j < len(action); j++ {
			var tmp []string
			tmp = append(tmp, resource[i])
			tmp = append(tmp, action[j])
			resp = append(resp, tmp)
		}
	}
	fmt.Println(resp)
	return http.StatusOK, resp, lg, false, nil
}

// view plan
func viewUserPlan(c echo.Context, request interface{}) (statusCode int, data interface{}, lg *model.LogFormat, logResponse bool, err error) {
	ctx := c.Request().Context()
	req := request.(*LoadUserUpdateReq)
	lg = &model.LogFormat{Source: c.Request().RemoteAddr, Action: "View User Plan", Data: req.ID}
	db := core.GetDB()
	uDAO := dao.GetUserDAO()

	done, er := uDAO.SelectUserDonePlan(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}

	onGoing, er := uDAO.SelectUserOnGoingPlan(ctx, db, req.ID)
	if er != nil {
		statusCode, err = http.StatusInternalServerError, er
		return
	}
	var result = &model.UserPlan{
		OnGoingPlan: onGoing,
		DonePlan: done,
	}
	
	return http.StatusOK, result, lg, false, nil
}

