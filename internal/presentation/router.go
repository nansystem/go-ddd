// HTTPルーティングの定義
// ハンドラーとのマッピング
package presentation

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	custommiddleware "github.com/nansystem/go-ddd/internal/presentation/middleware"
)

func NewRouter() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(custommiddleware.ErrorHandlerMiddleware())
	return e
}
