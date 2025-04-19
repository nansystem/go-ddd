package domainerror

import (
	"errors"
	"fmt"
)

// 基本的なエラー種類の定義
var (
	// ErrNotFound はエンティティが見つからない場合のエラーです
	ErrNotFound = errors.New("エンティティが見つかりません")

	// ErrInvalidInput は入力値が不正な場合のエラーです
	ErrInvalidInput = errors.New("入力値が不正です")

	// ErrDuplicated は重複エラーです
	ErrDuplicated = errors.New("既に存在します")

	// ErrUnauthorized は権限エラーです
	ErrUnauthorized = errors.New("権限がありません")

	// ErrInternal は内部エラーです
	ErrInternal = errors.New("内部エラーが発生しました")

	// データベース関連エラー
	ErrDatabase    = errors.New("データベースエラーが発生しました")
	ErrConnection  = errors.New("データベース接続エラーが発生しました")
	ErrTransaction = errors.New("トランザクションエラーが発生しました")
	ErrQuery       = errors.New("クエリ実行エラーが発生しました")
)

// NotFoundError は特定のエンティティが見つからないことを示すエラーです
type NotFoundError struct {
	EntityName string
	ID         string
}

// Error はエラーメッセージを返します
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s (ID: %s) %v", e.EntityName, e.ID, ErrNotFound)
}

// Is はエラー比較を行います
func (e *NotFoundError) Is(target error) bool {
	return target == ErrNotFound
}

// NewNotFoundError は新しいNotFoundErrorを作成します
func NewNotFoundError(entityName, id string) *NotFoundError {
	return &NotFoundError{
		EntityName: entityName,
		ID:         id,
	}
}

// ValidationError は入力値の検証エラーです
type ValidationError struct {
	Field   string
	Message string
}

// Error はエラーメッセージを返します
func (e *ValidationError) Error() string {
	return fmt.Sprintf("Field %s: %s", e.Field, e.Message)
}

// Is はエラー比較を行います
func (e *ValidationError) Is(target error) bool {
	return target == ErrInvalidInput
}

// NewValidationError は新しいValidationErrorを作成します
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// DuplicateEntryError は重複エントリのエラーを表します
type DuplicateEntryError struct {
	ID   string
	Name string
}

// Error はエラーメッセージを返します
func (e *DuplicateEntryError) Error() string {
	return fmt.Sprintf("重複エラー: ID=%s, Name=%s", e.ID, e.Name)
}

// Is はエラー比較を行います
func (e *DuplicateEntryError) Is(target error) bool {
	return target == ErrDuplicated
}

// NewDuplicateEntryError は新しいDuplicateEntryErrorを作成します
func NewDuplicateEntryError(id, name string) *DuplicateEntryError {
	return &DuplicateEntryError{
		ID:   id,
		Name: name,
	}
}

// DatabaseError はデータベース操作に関するエラーを表します
type DatabaseError struct {
	Operation string // 実行しようとした操作 (select, insert, update, delete など)
	Table     string // 対象テーブル
	Err       error  // 元のエラー
}

// Error はエラーメッセージを返します
func (e *DatabaseError) Error() string {
	return fmt.Sprintf("データベースエラー: 操作=%s, テーブル=%s, エラー=%v", e.Operation, e.Table, e.Err)
}

// Is はエラー比較を行います
func (e *DatabaseError) Is(target error) bool {
	return target == ErrDatabase || target == ErrConnection || target == ErrTransaction || target == ErrQuery
}
