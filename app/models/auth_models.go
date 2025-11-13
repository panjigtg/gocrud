package models

import(
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Users struct {
	ID        int       `json:"id" bson:"id"`
	Email     string    `json:"email" bson:"email"`
	Role      string    `json:"role" bson:"role"`
	Username  string    `json:"username" bson:"username"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}


type LoginRequest struct {
	UsernameOrEmail    	string 	`json:"1_username_or_email" validate:"required,email" example:"panji@example.com"`
	Password 			string 	`json:"2_password" validate:"required,min=8" example:"xafasfa122"`
}

type LoginResponse struct {
	User  Users   `json:"user"`
	Token string `json:"token"`
}


type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
 