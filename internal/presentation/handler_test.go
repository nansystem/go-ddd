package presentation_test

import (
	"bytes" // JSONEqのために必要
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/nansystem/go-ddd/internal/domain/domainerror"
	"github.com/nansystem/go-ddd/internal/domain/user"
	"github.com/nansystem/go-ddd/internal/presentation"
	"github.com/nansystem/go-ddd/internal/presentation/middleware" // エラーミドルウェアをインポート
	"github.com/nansystem/go-ddd/internal/usecase"
)

// --- テストヘルパー: Echoインスタンスとミドルウェア、ルートを設定 ---
func setupTestRouter(handler *presentation.UserHandler) *echo.Echo {
	e := echo.New()
	// エラーハンドリングミドルウェアを登録
	e.Use(middleware.ErrorHandlerMiddleware())

	// テスト対象のハンドラが処理するルートを登録
	g := e.Group("/users") // UserHandlerのSetupUserRoutesに合わせる
	handler.SetupUserRoutes(g)

	return e
}

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(mockService *usecase.MockUserService)
		expectedStatus int
		expectedBody   string // 期待するJSON文字列
	}{
		{
			name: "成功: ユーザー一覧を取得",
			setupMock: func(mockService *usecase.MockUserService) {
				users := []*user.User{
					{ID: "1", Name: "テストユーザー1"},
					{ID: "2", Name: "テストユーザー2"},
				}
				mockService.On("GetUsers").Return(users, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"ID":"1","Name":"テストユーザー1"},{"ID":"2","Name":"テストユーザー2"}]`,
		},
		{
			name: "失敗: ユースケースでエラー発生",
			setupMock: func(mockService *usecase.MockUserService) {
				// 内部エラーをシミュレート (DBエラーなど)
				mockService.On("GetUsers").Return(nil, errors.New("予期せぬ内部エラー")).Once()
			},
			expectedStatus: http.StatusInternalServerError, // ミドルウェアが500を返す
			expectedBody:   `{"error":"internal_server_error","message":"内部エラーが発生しました"}`,
		},
		{
			name: "成功: ユーザーが0件の場合",
			setupMock: func(mockService *usecase.MockUserService) {
				users := []*user.User{} // 空のスライス
				mockService.On("GetUsers").Return(users, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[]`, // 空のJSON配列
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(usecase.MockUserService)
			tt.setupMock(mockService)
			handler := presentation.NewUserHandler(mockService)
			e := setupTestRouter(handler)

			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(mockService *usecase.MockUserService, id string)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "成功: 存在するユーザーID",
			userID: "1",
			setupMock: func(mockService *usecase.MockUserService, id string) {
				user := &user.User{ID: id, Name: "テストユーザー1"}
				mockService.On("GetUserByID", id).Return(user, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ID":"1","Name":"テストユーザー1"}`,
		},
		{
			name:   "失敗: 存在しないユーザーID",
			userID: "notfound",
			setupMock: func(mockService *usecase.MockUserService, id string) {
				notFoundErr := domainerror.NewNotFoundError("User", id)
				mockService.On("GetUserByID", id).Return(nil, notFoundErr).Once()
			},
			expectedStatus: http.StatusNotFound, // ミドルウェアが404を返す
			expectedBody:   `{"error":"not_found","message":"User (ID: notfound) エンティティが見つかりません"}`,
		},
		{
			name:   "失敗: ユースケースで内部エラー発生",
			userID: "internalerror",
			setupMock: func(mockService *usecase.MockUserService, id string) {
				mockService.On("GetUserByID", id).Return(nil, errors.New("内部エラー発生")).Once()
			},
			expectedStatus: http.StatusInternalServerError, // ミドルウェアが500を返す
			expectedBody:   `{"error":"internal_server_error","message":"内部エラーが発生しました"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(usecase.MockUserService)
			tt.setupMock(mockService, tt.userID)
			handler := presentation.NewUserHandler(mockService)
			e := setupTestRouter(handler)

			targetURL := "/users/" + tt.userID
			req := httptest.NewRequest(http.MethodGet, targetURL, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		setupMock      func(mockService *usecase.MockUserService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "成功: ユーザーを作成",
			requestBody: `{"id":"newid","name":"新規ユーザー"}`,
			setupMock: func(mockService *usecase.MockUserService) {
				// CreateUserに渡されるであろうUserオブジェクトを期待値として設定
				expectedUser := &user.User{ID: "newid", Name: "新規ユーザー"}
				mockService.On("CreateUser", expectedUser).Return(nil).Once()
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"newid","message":"ユーザーが作成されました"}`, // handlerの実装に合わせる
		},
		{
			name:        "失敗: 不正なリクエストボディ (JSON)",
			requestBody: `{"id":"bad", "name":}`, // 不正なJSON
			setupMock: func(_ *usecase.MockUserService) {
				// Bindエラーなので、Usecaseは呼ばれない
			},
			expectedStatus: http.StatusBadRequest, // EchoのデフォルトのBindエラーは400
			// Echoのデフォルトエラーレスポンスか、ミドルウェアのレスポンスを期待
			// ここではミドルウェアが echo.ErrBadRequest を捕捉することを期待
			expectedBody: `{"error":"bad_request","message":"不正なリクエストです"}`,
		},
		{
			name:        "失敗: バリデーションエラー (Usecase)",
			requestBody: `{"id":"validid","name":""}`, // Nameが空
			setupMock: func(mockService *usecase.MockUserService) {
				invalidUser := &user.User{ID: "validid", Name: ""}
				validationErr := domainerror.NewValidationError("Name", "名前は必須です")
				mockService.On("CreateUser", invalidUser).Return(validationErr).Once()
			},
			expectedStatus: http.StatusBadRequest, // ミドルウェアが400を返す
			expectedBody:   `{"error":"invalid_input","message":"Field Name: 名前は必須です"}`,
		},
		{
			name:        "失敗: 重複エラー (Usecase)",
			requestBody: `{"id":"duplicateid","name":"重複ユーザー"}`,
			setupMock: func(mockService *usecase.MockUserService) {
				duplicateUser := &user.User{ID: "duplicateid", Name: "重複ユーザー"}
				duplicateErr := domainerror.NewDuplicateEntryError("duplicateid", "重複ユーザー")
				mockService.On("CreateUser", duplicateUser).Return(duplicateErr).Once()
			},
			expectedStatus: http.StatusConflict,                                                          // ミドルウェアが409を返す
			expectedBody:   `{"error":"duplicate_entry","message":"重複エラー: ID=duplicateid, Name=重複ユーザー"}`, // メッセージ調整
		},
		{
			name:        "失敗: その他の内部エラー (Usecase)",
			requestBody: `{"id":"internal","name":"内部エラー"}`,
			setupMock: func(mockService *usecase.MockUserService) {
				internalUser := &user.User{ID: "internal", Name: "内部エラー"}
				mockService.On("CreateUser", internalUser).Return(errors.New("予期せぬDBエラー")).Once()
			},
			expectedStatus: http.StatusInternalServerError, // ミドルウェアが500を返す
			expectedBody:   `{"error":"internal_server_error","message":"内部エラーが発生しました"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(usecase.MockUserService)
			tt.setupMock(mockService)
			handler := presentation.NewUserHandler(mockService)
			e := setupTestRouter(handler)

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON) // Content-Typeを設定
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}
