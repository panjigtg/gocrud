package services

import(
	"crudprojectgo/app/models"
	"crudprojectgo/app/repository"
	// "crudprojectgo/helper"
	"crudprojectgo/utils"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type AuthServices struct {
	repo *repository.AuthRepository
}

func NewAuthServices(repo *repository.AuthRepository) *AuthServices {
	return &AuthServices{
		repo: repo,
	}
}

func (s *AuthServices) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	user, passwordHash, err := s.repo.GetUserByUsernameOrEmail(req.UsernameOrEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("email atau password salah")
		}
		return &models.LoginResponse{}, err
	}

	if !utils.CheckPassword(req.Password, passwordHash) {
		return &models.LoginResponse{}, errors.New("Username atau password salah")
	}

	// Generate JWT
	token, err := utils.GenerateAccessToken(user)
	if err != nil {
		return &models.LoginResponse{}, errors.New("Gagal generate token")
	}

	return &models.LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *AuthServices) Register(req models.RegisterRequest) (models.Users, error) {
	// Cek apakah username/email sudah ada
	exist, _ := s.repo.CheckUserExist(req.Username, req.Email)
	if exist {
		return models.Users{}, errors.New("Username atau email sudah digunakan")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return models.Users{}, errors.New("Gagal mengenkripsi password")
	}

	// Set default role if empty
	role := req.Role
	if role == "" {
		role = "user"
	}

	newUser, err := s.repo.CreateUserWithRole(req.Username, req.Email, hashedPassword, role)
	if err != nil {
		return models.Users{}, errors.New("Gagal membuat user")
	}
	return newUser, nil
}

func (s *AuthServices) RefreshAccessToken(refreshToken string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.GetJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("refresh token tidak valid atau expired")
	}

	// Ambil user dari DB pakai claims.Subject (username)
	user, _, err := s.repo.GetUserByUsernameOrEmail(claims.Subject)
	if err != nil {
		return "", errors.New("pengguna tidak ditemukan")
	}

	return utils.GenerateAccessToken(user)
}
