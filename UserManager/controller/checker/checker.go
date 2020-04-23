package checker

import (
	"DBproject1/controller"

	"github.com/labstack/echo"
)

type reqLoginCheck struct {
	Username string `json:"username"` // omitempty
	Password string `json:"password"`
}

type reqPermissionCheck struct {
	Token  string `json:"token"`
	Des    string `json:"des"`
	Action string `json:"action"`
}

func CheckLogin(c echo.Context) error {
	return controller.ExecHandler(c, &reqLoginCheck{}, checkLogin)
}

func CheckPermission(c echo.Context) error {
	return controller.ExecHandler(c, &reqPermissionCheck{}, checkPermission)
}
