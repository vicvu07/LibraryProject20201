package resource

import (
	"DBproject1/controller"

	"github.com/labstack/echo"
)

type viewRequest struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

type resource struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

func Save(e echo.Context) error {
	return controller.ExecHandler(e, &resource{}, save)
}

func Delete(e echo.Context) error {
	return controller.ExecHandler(e, &resource{}, delete)
}

func SelectByView(e echo.Context) error {
	return controller.ExecHandler(e, &viewRequest{}, selectByView)
}

func SelectAll(e echo.Context) error {
	return controller.ExecHandler(e, nil, selectAll)
}
