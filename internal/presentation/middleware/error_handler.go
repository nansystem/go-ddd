package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/nansystem/go-ddd/internal/domain/domainerror"
)

// ErrorResponse はエラーレスポンスの形式を定義します
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// ErrorHandlerMiddleware はAPIエラーハンドリングのミドルウェアです
func ErrorHandlerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 次のハンドラを実行
			err := next(c)

			// エラーがなければ処理終了
			if err == nil {
				return nil
			}

			// エラーの種類に応じてHTTPステータスとメッセージを設定
			var statusCode int
			var response ErrorResponse

			// 独自のエラータイプを判別
			var notFoundErr *domainerror.NotFoundError
			var duplicateErr *domainerror.DuplicateEntryError
			var validationErr *domainerror.ValidationError

			// エラータイプに基づいてレスポンスを構築
			switch {
			case errors.Is(err, domainerror.ErrNotFound) || errors.As(err, &notFoundErr):
				statusCode = http.StatusNotFound
				response.Error = "not_found"
				response.Message = err.Error()

			case errors.Is(err, domainerror.ErrDuplicated) || errors.As(err, &duplicateErr):
				statusCode = http.StatusConflict
				response.Error = "duplicate_entry"
				response.Message = err.Error()

			case errors.Is(err, domainerror.ErrInvalidInput) || errors.As(err, &validationErr):
				statusCode = http.StatusBadRequest
				response.Error = "invalid_input"
				response.Message = err.Error()

			case errors.Is(err, domainerror.ErrUnauthorized):
				statusCode = http.StatusUnauthorized
				response.Error = "unauthorized"
				response.Message = err.Error()

			// データベース関連エラーは内部エラーとして扱う
			case errors.Is(err, domainerror.ErrDatabase) ||
				errors.Is(err, domainerror.ErrConnection) ||
				errors.Is(err, domainerror.ErrTransaction) ||
				errors.Is(err, domainerror.ErrQuery):
				statusCode = http.StatusInternalServerError
				response.Error = "internal_server_error"
				response.Message = "内部エラーが発生しました"

				// ログにエラー詳細を記録（本番環境ではクライアントに詳細を返さない）
				c.Logger().Error(err)

			case errors.Is(err, echo.ErrNotFound):
				statusCode = http.StatusNotFound
				response.Error = "not_found"
				response.Message = "リソースが見つかりません"

			case errors.Is(err, echo.ErrMethodNotAllowed):
				statusCode = http.StatusMethodNotAllowed
				response.Error = "method_not_allowed"
				response.Message = "許可されていないメソッドです"

			case errors.Is(err, echo.ErrBadRequest):
				statusCode = http.StatusBadRequest
				response.Error = "bad_request"
				response.Message = "不正なリクエストです"

			default:
				// その他のエラー
				statusCode = http.StatusInternalServerError
				response.Error = "internal_server_error"
				response.Message = "内部エラーが発生しました"

				// 開発環境では元のエラーメッセージも含める
				if c.Echo().Debug {
					response.Message = err.Error()
				}

				// ログにエラー詳細を記録
				c.Logger().Error(err)
			}

			// JSONレスポンスを返す
			return c.JSON(statusCode, response)
		}
	}
}
