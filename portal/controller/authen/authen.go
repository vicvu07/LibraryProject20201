package authen

import (
	"github.com/labstack/echo"
	"github.com/pinezapple/LibraryProject20201/portal/controller"
)

type reqLogin struct {
	Username string `json:"username"` // omitempty
	Password string `json:"password"`
}

// Login do login for user
func Login(c echo.Context) (erro error) {
	return controller.ExecHandler(c, &reqLogin{}, login)
}
