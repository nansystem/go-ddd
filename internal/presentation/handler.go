// エンドポイント固有のハンドラー
// コントローラーの役割
package presentation

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	return c.String(http.StatusOK, "GetUsers response")
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
