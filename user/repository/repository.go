package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"imp/assessment/user/entity"
)

type userRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) List(ctx context.Context) (res []*entity.User, err error) {
	start := "SELECT id, username, password, created_at, updated_at, deleted_at FROM users WHERE deleted_at IS NULL"

	var sql strings.Builder
	sql.WriteString(start)

	// add limit and offset
	sql.WriteString(fmt.Sprintf(` LIMIT %d OFFSET %d`, 10, 0))
	rows, err := r.db.Query(sql.String())
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt); err != nil {
			return nil, err
		}

		res = append(res, &user)
	}

	return res, nil
}

func (r *userRepository) FindOneByUsername(ctx context.Context, username string) (res *entity.User, err error) {
	sql := "SELECT id, username, fullname, password, created_at, updated_at, deleted_at FROM users WHERE username = $1 AND deleted_at IS NULL"

	row := r.db.QueryRow(sql, username)

	user := entity.User{}

	if err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Fullname,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt); err != nil {

		return nil, nil // no record
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, m *entity.User) (res *entity.User, err error) {
	sql := `INSERT INTO users (id, username, fullname, password)
			VALUES ($1, $2, $3, $4)
			RETURNING id, username, fullname, password, created_at, updated_at, deleted_at`

	row := r.db.QueryRow(sql, m.Id, m.Username, m.Fullname, m.Password)

	if err := row.Scan(
		&m.Id,
		&m.Username,
		&m.Fullname,
		&m.Password,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.DeletedAt); err != nil {
		return nil, err
	}

	return m, nil
}
