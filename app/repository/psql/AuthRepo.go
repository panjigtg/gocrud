package psql

import (
	"crudprojectgo/app/models"
	"database/sql"
	"time"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

// ✅ Method untuk login (username/email)
func (r *AuthRepository) GetUserByUsernameOrEmail(identifier string) (models.Users, string, error) {
	var user models.Users
	var passwordHash string

	err := r.db.QueryRow(`
		SELECT id, username, email, password_hash, role, created_at 
		FROM users 
		WHERE username = $1 OR email = $1 LIMIT 1
	`, identifier).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return user, "", err
	}

	return user, passwordHash, nil
}

func (r *AuthRepository) CheckUserExist(username, email string) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM users WHERE username = $1 OR email = $2
	`, username, email).Scan(&count)

	return count > 0, err
}


func (r *AuthRepository) CreateUser(username, email, passwordHash string) (models.Users, error) {
	var user models.Users

	err := r.db.QueryRow(`
		INSERT INTO users (username, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, 'user', $4)
		RETURNING id, username, email, role, created_at
	`, username, email, passwordHash, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *AuthRepository) CreateUserWithRole(username, email, passwordHash, role string) (models.Users, error) {
	var user models.Users

	err := r.db.QueryRow(`
		INSERT INTO users (username, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, username, email, role, created_at
	`, username, email, passwordHash, role).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
	)

	return user, err
}



