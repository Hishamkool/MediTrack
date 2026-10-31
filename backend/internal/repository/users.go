package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"meditrack/internal/models"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

const userColumns = `id, name, email, password_hash, role, is_active, created_at`

// scanUser copies one row into a models.User.
func scanUser(row pgx.Row) (*models.User, error) {
	var u models.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	row := r.db.QueryRow(ctx, `SELECT `+userColumns+` FROM users WHERE id = $1`, id)
	return scanUser(row)
}

// GetByEmail expects an already lower-cased email.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.QueryRow(ctx, `SELECT `+userColumns+` FROM users WHERE email = $1`, email)
	return scanUser(row)
}

func (r *UserRepo) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.Query(ctx, `SELECT `+userColumns+` FROM users ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{} // empty slice -> JSON [] instead of null
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, rows.Err()
}

// Create inserts a user and fills in ID and CreatedAt.
func (r *UserRepo) Create(ctx context.Context, u *models.User) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (name, email, password_hash, role, is_active)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at`,
		u.Name, u.Email, u.PasswordHash, u.Role, u.IsActive,
	).Scan(&u.ID, &u.CreatedAt)

	// 23505 is Postgres' "unique_violation" error code (duplicate email).
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrDuplicate
	}
	return err
}

// Update saves name, role, is_active and password_hash for an existing user.
func (r *UserRepo) Update(ctx context.Context, u *models.User) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE users SET name = $1, role = $2, is_active = $3, password_hash = $4
		 WHERE id = $5`,
		u.Name, u.Role, u.IsActive, u.PasswordHash, u.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// CountActiveAdmins is used to avoid locking everyone out, and to decide
// whether the first admin needs to be seeded.
func (r *UserRepo) CountActiveAdmins(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRow(ctx,
		`SELECT count(*) FROM users WHERE role = 'admin' AND is_active`,
	).Scan(&n)
	return n, err
}
