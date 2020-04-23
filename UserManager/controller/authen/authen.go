package authen

import (
	"DBproject1/controller"

	"github.com/labstack/echo"
)

type reqLogin struct {
	Username string `json:"username"` // omitempty
	Password string `json:"password"`
}

// Login do login for user
func Login(c echo.Context) (erro error) {
	return controller.ExecHandler(c, &reqLogin{}, login)
}
