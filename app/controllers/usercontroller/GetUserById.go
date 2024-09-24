package usercontroller

import (
    "github.com/labstack/echo/v4"
)

type GetUserByIdRespone struct {

}

// @method		Get
// @path		/User/:id
// @response	GetUserByIdRespone
func (u *UserController) GetUserById(ctx echo.Context) error {
    return nil
}