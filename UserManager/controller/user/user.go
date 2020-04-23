package user

import (
	"DBproject1/controller"

	"github.com/labstack/echo"
)

type UserDeleteReq struct {
	ID       uint64 `json:"ID"`
	Username string `json:"Username"` // omitempty
}

type UpdatePermissionReq struct {
	Username   string     `json:"Username"` // omitempty
	Permission [][]string `json:"Permission"`
}

type userSaveReq struct {
	Name         string `json:"Name"`
	DOB          string `json:"DOB"`
	Sex          string `json:"Sex"`
	Position     string `json:"Position"`
	PhoneNum     string `json:"Phonenum"`
	NationalID   string `json:"NationalID"`
	Salary       string `json:"Salary"`
	Username     string `json:"Username"`
	DepartmentID uint64 `json:"DepartmentID"`
	Password     string `json:"Password"`
}

type userUpdateReq struct {
	ID         uint64 `json:"ID"`
	Name       string `json:"Name"`
	DOB        string `json:"DOB"`
	Sex        string `json:"Sex"`
	Position   string `json:"Position"`
	PhoneNum   string `json:"Phonenum"`
	NationalID string `json:"NationalID"`
	Salary     string `json:"Salary"`
}

type UpdateRoleReq struct {
	Username string   `json:"Username"` // omitempty
	Role     []string `json:"Role"`
}

type LoadUserUpdateReq struct {
	ID uint64 `json:"ID"`
}

type UpdateUserDepartmentReq struct {
	ID              uint64 `json:"ID"`
	OldDepartmentID uint64 `json:"OldDepartmentID"`
	NewDepartmentID uint64 `json:"NewDepartmentID"`
}

type viewRequest struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

type loadUpdateFormReq struct {
	Username string `json:"username"`
}

type reqChangePasswd struct {
	Username    string `json:"Username"`
	OldPassword string `json:"OldPassword"`
	NewPassword string `json:"NewPassword"`
}

func DeleteUser(e echo.Context) error {
	return controller.ExecHandler(e, &UserDeleteReq{}, deleteUser)
}

func Save(e echo.Context) error {
	return controller.ExecHandler(e, &userSaveReq{}, save)
}

func ChangePassword(e echo.Context) error {
	return controller.ExecHandler(e, &reqChangePasswd{}, changePassword)
}

func SelectUserByView(e echo.Context) error {
	return controller.ExecHandler(e, &viewRequest{}, selectUserByView)
}

func LoadUpdatePermissionForm(e echo.Context) error {
	return controller.ExecHandler(e, &loadUpdateFormReq{}, loadUpdatePermissionForm)
}

func LoadUpdateRoleForm(e echo.Context) error {
	return controller.ExecHandler(e, &loadUpdateFormReq{}, loadUpdateRoleForm)
}

func LoadUpdateUserDetailForm(e echo.Context) error {
	return controller.ExecHandler(e, &LoadUserUpdateReq{}, loadUpdateUserDetailForm)
}

func LoadUpdateUserDepartmentForm(e echo.Context) error {
	return controller.ExecHandler(e, &LoadUserUpdateReq{}, loadUpdateUserDeparmentForm)
}

func ViewUserPlan(e echo.Context) error {
	return controller.ExecHandler(e, &LoadUserUpdateReq{}, viewUserPlan)
}

func UpdatePermission(e echo.Context) error {
	return controller.ExecHandler(e, &UpdatePermissionReq{}, updatePermission)
}

func SelectAllUser(e echo.Context) error {
	return controller.ExecHandler(e, nil, selectAllUser)
}

func UpdateRole(e echo.Context) error {
	return controller.ExecHandler(e, &UpdateRoleReq{}, updateRole)
}

func UpdateUserDepartment(e echo.Context) error {
	return controller.ExecHandler(e, &UpdateUserDepartmentReq{}, updateUserDepartment)
}

func UpdateUserDetail(e echo.Context) error {
	return controller.ExecHandler(e, &userUpdateReq{}, updateUserDetail)
}

func SelectAllPermission(e echo.Context) error {
	return controller.ExecHandler(e, nil, selectAllPermission)
}
