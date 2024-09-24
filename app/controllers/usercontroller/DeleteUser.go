package usercontroller

import (
    "github.com/labstack/echo/v4"
)

type DeleeteUserRespone struct {

}

// @method		DELETE
// @path		/User/:id
// @response	DeleeteUserRespone
func (u *UserController) DeleteUser (ctx echo.Context) error {
    return nil
}