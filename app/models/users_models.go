package models

import(
	"time"
)

type FilterUsers struct {
	ID 			int 		`json:"id"`
	Username 	string 		`json:"username"`
	Email 		string 		`json:"email"`
	CreatedAt 	time.Time 	`json:"created_at"`
}

