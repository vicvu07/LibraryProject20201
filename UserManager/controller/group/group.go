package group

import (
	"DBproject1/controller"

	"github.com/labstack/echo"
)

type reqGroup struct {
	GroupName   string `json:"GroupName"`
	Description string `json:"Description"`
}

type reqLoadGroupPermission struct {
	GroupName string `json:"GroupName"`
}
type reqGroupPermission struct {
	GroupName  string     `json:"GroupName"`
	Permission [][]string `json:"Permission"`
}

type viewRequest struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

func CreateGroup(e echo.Context) error {
	return controller.ExecHandler(e, &reqGroup{}, createGroup)
}

func DeleteGroup(e echo.Context) error {
	return controller.ExecHandler(e, &reqGroup{}, deleteGroup)
}

func UpdateGroupPermission(e echo.Context) error {
	return controller.ExecHandler(e, &reqGroupPermission{}, updateGroupPermission)
}

func SelectGroupByView(e echo.Context) error {
	return controller.ExecHandler(e, &viewRequest{}, selectGroupByView)
}

func SelectAllGroup(e echo.Context) error {
	return controller.ExecHandler(e, nil, selectAllGroup)
}

func LoadGroupPermission(e echo.Context) error {
	return controller.ExecHandler(e, &reqLoadGroupPermission{}, loadGroupPermission)
}
