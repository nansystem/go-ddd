package presentation_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/nansystem/go-ddd/internal/domain/domainerror"
	"github.com/nansystem/go-ddd/internal/domain/user"
	"github.com/nansystem/go-ddd/internal/presentation"
	"github.com/nansystem/go-ddd/internal/usecase"
)

func TestGetUsers(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		setupMock      func(mock *usecase.MockUserService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "ユーザー一覧の取得に成功",
			setupMock: func(mockService *usecase.MockUserService) {
				users := []*user.User{
					{ID: "1", Name: "テストユーザー1"},
					{ID: "2", Name: "テストユーザー2"},
				}
				mockService.On("GetUsers").Return(users, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"ID":"1","Name":"テストユーザー1"},{"ID":"2","Name":"テストユーザー2"}]`,
		},
		{
			name: "エラーが発生した場合",
			setupMock: func(mockService *usecase.MockUserService) {
				mockService.On("GetUsers").Return(nil, errors.New("データベースエラー"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"データベースエラー"`,
		},
	}

	// テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Echo インスタンスとリクエストを設定
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// モックサービスを設定
			mockService := new(usecase.MockUserService)
			tt.setupMock(mockService)

			// ハンドラーを作成してテスト
			handler := presentation.NewUserHandler(mockService)
			err := handler.GetUsers(c)

			// アサーション
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// JSONレスポンスを整形して比較（スペースや改行の違いを無視）
			var expected, actual interface{}
			_ = json.Unmarshal([]byte(tt.expectedBody), &expected)
			_ = json.Unmarshal(rec.Body.Bytes(), &actual)
			assert.Equal(t, expected, actual)

			// モックの呼び出しを検証
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupMock      func(mock *usecase.MockUserService, id string)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "存在するユーザーのIDを指定",
			userID: "1",
			setupMock: func(mockService *usecase.MockUserService, id string) {
				user := &user.User{ID: id, Name: "テストユーザー"}
				mockService.On("GetUserByID", id).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ID":"1","Name":"テストユーザー"}`,
		},
		{
			name:   "存在しないユーザーのIDを指定",
			userID: "999",
			setupMock: func(mockService *usecase.MockUserService, id string) {
				mockService.On("GetUserByID", id).Return(nil, domainerror.ErrNotFound)
			},
			expectedStatus: http.StatusInternalServerError, // エラー時は500を返すように修正
			expectedBody:   `"ユーザーが見つかりません"`,               // domainerror.ErrNotFound のメッセージ
		},
		{
			name:   "GetUserByIDで予期せぬエラーが発生",
			userID: "err",
			setupMock: func(mockService *usecase.MockUserService, id string) {
				mockService.On("GetUserByID", id).Return(nil, errors.New("予期せぬDBエラー"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"予期せぬDBエラー"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Echo インスタンスとリクエストを設定
			e := echo.New()
			// パスパラメータを含むリクエストURLを正しく設定
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.userID)

			// モックサービスを設定
			mockService := new(usecase.MockUserService)
			tt.setupMock(mockService, tt.userID)

			// ハンドラーを作成してテスト
			handler := presentation.NewUserHandler(mockService)
			err := handler.GetUserByID(c)

			// アサーション
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// レスポンスボディを確認
			assert.Equal(t, tt.expectedBody, strings.TrimSpace(rec.Body.String()))
		})
	}
}
