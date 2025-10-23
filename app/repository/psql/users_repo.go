package psql

import (
	"crudprojectgo/app/models"
	"database/sql"
	"fmt"
	"log"
	"strings"

	// "crudprojectgo/helper"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
    return &UsersRepository{db: db}
}


func (r *UsersRepository) GetUsersRepo(search string, sortBy, order string, page, limit, offset int) ([]models.FilterUsers, error) {
	validSort := map[string]string{
		"id": "id", "username": "username", "email": "email", "created_at": "created_at",
	}
	sortCol, ok := validSort[sortBy]
	if !ok {
		sortCol = "id"
	}

	ord := "ASC"
	if strings.ToUpper(order) == "DESC" {
		ord = "DESC"
	}

	query := fmt.Sprintf(`
		SELECT id, username, email, created_at
		FROM users 
		WHERE username ILIKE $1 OR email ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortCol, ord)

	rows, err := r.db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Println("query Error:", err)
		return nil, err
	}

	defer rows.Close()

	var filterUsers []models.FilterUsers	
	for rows.Next() {
		var f models.FilterUsers
		if err := rows.Scan(&f.ID, &f.Username, &f.Email, &f.CreatedAt);
		err != nil {
			return nil, err
		}
		filterUsers = append(filterUsers, f)
	}
	return filterUsers, nil
}

func (r *UsersRepository) CountUsersRepo(search string) (int, error) {
	var total int
	countQuery := `
		Select Count(*) FROM users WHERE username ILIKE $1 OR email ILIKE $1
	`
	err := r.db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
	return 0, err
	}
	return total, nil
}