package models

import(
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Users struct {
	ID 	  int       `json:"id"`
	Email string    `json:"email"`
	// Password string  `json:"password"`
	Role    string    `json:"role"`
	Username string  `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt time.Time `json:"updated_at"`	
}

type LoginRequest struct {
	UsernameOrEmail    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
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


// type 