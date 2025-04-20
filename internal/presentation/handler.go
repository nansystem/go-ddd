package presentation

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/nansystem/go-ddd/internal/domain/user"
	"github.com/nansystem/go-ddd/internal/usecase"
)

type UserHandler struct {
	userService usecase.UserServiceInterface
}

func NewUserHandler(userService usecase.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.userService.GetUsers()
	if err != nil {
		return err // エラーをそのまま返す
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		return err // エラーをそのまま返す
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	// リクエストボディのバインディングはハンドラで行うのが一般的
	reqUser := new(struct { // DTOを定義する方が望ましい場合もある
		ID    string `json:"id"` // Create時はIDは不要か、自動生成するべき
		Name  string `json:"name"`
		Email string `json:"email"` // UserドメインモデルにEmailがなければ不要
	})
	if err := c.Bind(reqUser); err != nil {
		// バインドエラーはBadRequestとしてミドルウェアに処理させるか、
		// ここで具体的なエラーを返す（ただし、ミドルウェアの思想とは少しずれる）
		// return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body: "+err.Error())
		return err // シンプルにミドルウェアに任せる
	}

	// ドメインモデルに変換 (Emailはドメインモデルに合わせて調整)
	domainUser := &user.User{
		ID:   reqUser.ID, // IDの扱いは要検討
		Name: reqUser.Name,
		// Email: reqUser.Email, // 必要なら追加
	}

	if err := h.userService.CreateUser(domainUser); err != nil {
		return err // エラーをそのまま返す
	}

	// 成功レスポンス (message フィールドを追加)
	return c.JSON(http.StatusCreated, map[string]string{
		"id":      domainUser.ID,
		"message": "ユーザーが作成されました",
	})
}

func (h *UserHandler) SetupUserRoutes(g *echo.Group) {
	g.GET("", h.GetUsers)
	g.GET("/:id", h.GetUserByID)
	g.POST("", h.CreateUser)
}
