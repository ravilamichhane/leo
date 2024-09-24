package usercontroller

import (
	"github.com/labstack/echo/v4"
)

type CreateUserRequest struct {
	Name string `formtype:"string"`
	Age  int    `formtype:"number"`
}

type CreateUserRespone struct {
}

// @method		POST
// @path		/User/:id
// @response	CreateUserRespone
// @body        CreateUserRequest
// @generateform
func (u *UserController) CreateUser(ctx echo.Context) error {
	return nil
}
