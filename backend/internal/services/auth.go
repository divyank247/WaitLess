package services

import (
	"database/sql"
	"fmt"
	"waitless-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db : db}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.User,error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w",err)
	}
	if exists {
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password),bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w",err)
	}

	if req.Role == "" {
		req.Role = "user"
	}

	user := &models.User{
		Email: req.Email,
		Password: string(hashedPassword),
		Name: req.Name,
		Role: req.Role,
	}

	query := `INSERT INTO users (email,password,name,role)
	VALUES ($1,$2,$3,$4)
	RETURNING id,created_at`

	err = s.db.QueryRow(query,user.Email,user.Password,user.Name,user.Role).Scan(&user.ID,&user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w",err)
	}

	return user, nil
}