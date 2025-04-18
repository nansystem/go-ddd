package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func main() {
	// Echoインスタンスを作成
	e := echo.New()

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ルートハンドラ
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo!")
	})

	// GETエンドポイント
	e.GET("/users", getUsers)
	e.GET("/users/:id", getUserByID)

	// POSTエンドポイント
	e.POST("/users", createUser)

	// サーバー起動
	e.Logger.Fatal(e.Start(":8080"))
}

// GETエンドポイント実装
func getUsers(c echo.Context) error {
	users := []User{
		{Name: "山田太郎", Email: "taro@example.com"},
		{Name: "佐藤花子", Email: "hanako@example.com"},
	}

	return c.JSON(http.StatusOK, Response{
		Message: "ユーザー一覧取得成功",
		Data:    users,
	})
}

func getUserByID(c echo.Context) error {
	id := c.Param("id")

	// 簡易的なIDチェック（実際はDBなどで検索）
	if id == "1" {
		user := User{Name: "山田太郎", Email: "taro@example.com"}
		return c.JSON(http.StatusOK, Response{
			Message: "ユーザー取得成功",
			Data:    user,
		})
	}

	return c.JSON(http.StatusNotFound, Response{
		Message: "ユーザーが見つかりません",
	})
}

// POSTエンドポイント実装
func createUser(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Message: "リクエスト形式が不正です",
		})
	}

	// バリデーション
	if user.Name == "" || user.Email == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Message: "名前とメールアドレスは必須です",
		})
	}

	// 本来はDBに保存などの処理

	return c.JSON(http.StatusCreated, Response{
		Message: "ユーザー作成成功",
		Data:    user,
	})
}
