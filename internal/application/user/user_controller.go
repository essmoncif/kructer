package user

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"kructer.com/internal/context"
	"kructer.com/internal/core/errors"
)

type UserController struct {
	userService *UserService
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		userService: NewUserService(db),
	}
}

func (ctrl *UserController) CreateUser(c echo.Context) error {
	cc := c.(*context.KructerContext)

	user := new(User)
	if err := cc.Bind(user); err != nil {
		b := errors.NewBoom(errors.InvalidBindingModel, errors.ErrorText(errors.InvalidBindingModel), err)
		cc.Logger().Error(err)
		return cc.JSON(http.StatusBadRequest, b)
	}
	newUser, err := ctrl.userService.CreateUser(*user)
	if err != nil {
		b := errors.NewBoom(errors.EntityCreationError, errors.ErrorText(errors.EntityCreationError), err)
		cc.Logger().Error(err)
		return cc.JSON(http.StatusInternalServerError, b)
	}

	return cc.JSON(http.StatusCreated, newUser)
}