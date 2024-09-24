package usercontroller

import (
    "github.com/labstack/echo/v4"
)

type GetAllUsersRespone struct {

}

// @method		Get
// @path		/User
// @response	Page<GetAllUsersRespone>
// @generatedatatable
func (u *UserController) GetAllUsers(ctx echo.Context) error {
    return nil
}