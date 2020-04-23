package action

import (
	"DBproject1/controller"

	"github.com/labstack/echo"
)

type viewRequest struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

type action struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

func Save(e echo.Context) error {
	return controller.ExecHandler(e, &action{}, save)
}

func Delete(e echo.Context) error {
	return controller.ExecHandler(e, &action{}, delete)
}

func SelectByView(e echo.Context) error {
	return controller.ExecHandler(e, &viewRequest{}, selectByView)
}

func SelectAll(e echo.Context) error {
	return controller.ExecHandler(e, nil, selectAll)
}
