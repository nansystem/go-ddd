// エンドポイント固有のハンドラー
// コントローラーの役割
package presentation

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/nansystem/go-ddd/internal/usecase"
)

type UserHandler struct {
	userService *usecase.UserService
}

func NewUserHandler(userService *usecase.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.userService.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, "GetUserByID response for ID: "+id)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	return c.String(http.StatusCreated, "CreateUser response")
}

func (h *UserHandler) SetupUserRoutes(g *echo.Group) {
	g.GET("", h.GetUsers)
	g.GET("/:id", h.GetUserByID)
	g.POST("", h.CreateUser)
}
