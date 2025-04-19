package domainerror

import (
	"errors"
	"fmt"
)

// データベース関連のエラー定義
var (
	// ErrDatabase はデータベース操作に関する一般的なエラーです
	ErrDatabase = errors.New("データベースエラーが発生しました")

	// ErrConnection はデータベース接続エラーです
	ErrConnection = errors.New("データベース接続エラーが発生しました")

	// ErrTransaction はトランザクションエラーです
	ErrTransaction = errors.New("トランザクションエラーが発生しました")

	// ErrQuery はクエリ実行エラーです
	ErrQuery = errors.New("クエリ実行エラーが発生しました")

	// ErrDuplicateEntry は重複エラーです
	ErrDuplicateEntry = errors.New("重複エラーが発生しました")
)

// DatabaseError はデータベース操作関連のエラーを表します
type DatabaseError struct {
	Operation string // 実行しようとした操作 (select, insert, update, delete など)
	Table     string // 対象テーブル
	Err       error  // 元のエラー
}

// Error はエラーメッセージを返します
func (e *DatabaseError) Error() string {
	return fmt.Sprintf("データベースエラー: 操作=%s, テーブル=%s, 詳細=%v", e.Operation, e.Table, e.Err)
}

// Unwrap は元のエラーを返します（Go 1.13以降のエラーラッピング）
func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// Is はエラー比較を行います
func (e *DatabaseError) Is(target error) bool {
	return target == ErrDatabase
}

// NewDatabaseError は新しいDatabaseErrorを作成します
func NewDatabaseError(operation, table string, err error) *DatabaseError {
	return &DatabaseError{
		Operation: operation,
		Table:     table,
		Err:       err,
	}
}

// ConnectionError はデータベース接続エラーを表します
type ConnectionError struct {
	DSN string // 接続文字列（機密情報は含まない）
	Err error  // 元のエラー
}

// Error はエラーメッセージを返します
func (e *ConnectionError) Error() string {
	return fmt.Sprintf("データベース接続エラー: %v", e.Err)
}

// Unwrap は元のエラーを返します
func (e *ConnectionError) Unwrap() error {
	return e.Err
}

// Is はエラー比較を行います
func (e *ConnectionError) Is(target error) bool {
	return target == ErrConnection
}

// NewConnectionError は新しいConnectionErrorを作成します
func NewConnectionError(dsn string, err error) *ConnectionError {
	// DSNから機密情報（パスワードなど）を削除すべきですが、簡略化のため省略
	return &ConnectionError{
		DSN: dsn,
		Err: err,
	}
}

// QueryError はクエリ実行エラーを表します
type QueryError struct {
	Query string // 実行しようとしたクエリ（プリペアドステートメントの場合はプレースホルダーあり）
	Args  []any  // クエリパラメータ
	Err   error  // 元のエラー
}

// Error はエラーメッセージを返します
func (e *QueryError) Error() string {
	return fmt.Sprintf("クエリ実行エラー: %v", e.Err)
}

// Unwrap は元のエラーを返します
func (e *QueryError) Unwrap() error {
	return e.Err
}

// Is はエラー比較を行います
func (e *QueryError) Is(target error) bool {
	return target == ErrQuery
}

// NewQueryError は新しいQueryErrorを作成します
func NewQueryError(query string, args []any, err error) *QueryError {
	return &QueryError{
		Query: query,
		Args:  args,
		Err:   err,
	}
}

// TransactionError はトランザクションエラーを表します
type TransactionError struct {
	Operation string // トランザクション操作（begin, commit, rollback）
	Err       error  // 元のエラー
}

// Error はエラーメッセージを返します
func (e *TransactionError) Error() string {
	return fmt.Sprintf("トランザクションエラー (%s): %v", e.Operation, e.Err)
}

// Unwrap は元のエラーを返します
func (e *TransactionError) Unwrap() error {
	return e.Err
}

// Is はエラー比較を行います
func (e *TransactionError) Is(target error) bool {
	return target == ErrTransaction
}

// NewTransactionError は新しいTransactionErrorを作成します
func NewTransactionError(operation string, err error) *TransactionError {
	return &TransactionError{
		Operation: operation,
		Err:       err,
	}
}

type DuplicateEntryError struct {
	ID   string
	Name string
}

// Error はエラーメッセージを返します
func (e *DuplicateEntryError) Error() string {
	return "重複エラー: ID=" + e.ID + ", Name=" + e.Name
}

// Is はエラー比較を行います
func (e *DuplicateEntryError) Is(target error) bool {
	return target == ErrDuplicateEntry
}

// NewDuplicateEntryError は新しいDuplicateEntryErrorを作成します
func NewDuplicateEntryError(id, name string) *DuplicateEntryError {
	return &DuplicateEntryError{
		ID:   id,
		Name: name,
	}
}
