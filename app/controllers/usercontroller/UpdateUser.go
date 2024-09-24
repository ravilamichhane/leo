package usercontroller

import (
    "github.com/labstack/echo/v4"
)

type UpdateUserRequest struct {

}

type UpdateUserRespone struct {

}

// @method		PUT
// @path		/User/:id
// @response	UpdateUserRespone
// @body        UpdateUserRequest
// @generateform
func (u *UserController) UpdateUser (ctx echo.Context) error {
    return nil
}