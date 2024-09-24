package usercontroller

import (
    "github.com/labstack/echo/v4"
)

type UserController struct {
    // Inject Here
}

func New() *UserController {
    return &UserController{}
}

func (u *UserController) addRoutes(e echo.Group) {
        g := e.Group("/user")
        g.GET("/", u.GetAllUsers)
		g.GET("/:id", u.GetUserById)
		g.POST("", u.CreateUser)
		g.PUT("/:id", u.UpdateUser)
		g.DELETE("/:id", u.DeleteUser)
}