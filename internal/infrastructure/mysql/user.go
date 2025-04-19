package mysql

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"

	"github.com/nansystem/go-ddd/internal/domain/domainerror"
	"github.com/nansystem/go-ddd/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUsers() ([]*user.User, error) {
	rows, err := r.db.Query("SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}

	users := []*user.User{}
	for rows.Next() {
		var id string
		var name string
		err = rows.Scan(&id, &name)

		if err != nil {
			return nil, err
		}

		users = append(users, user.NewUser(id, name))
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(id string) (*user.User, error) {
	row := r.db.QueryRow("SELECT id, name FROM users WHERE id = $1", id)
	var name string
	err := row.Scan(&id, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainerror.ErrNotFound
		}
		return nil, err
	}

	return user.NewUser(id, name), nil
}

func (r *UserRepository) CreateUser(user *user.User) error {
	_, err := r.db.Exec("INSERT INTO users (id, name) VALUES ($1, $2)", user.ID, user.Name)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return domainerror.NewDuplicateEntryError(user.ID, user.Name)
			}
		}
		return err
	}
	return nil
}
