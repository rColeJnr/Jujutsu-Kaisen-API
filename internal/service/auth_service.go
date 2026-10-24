package service

import (
	"fmt"
	"strings"

	"github.com/rcolejnr/jjk-api/internal/jwt"
	"github.com/rcolejnr/jjk-api/internal/models"
	"github.com/rcolejnr/jjk-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(req *models.LoginRequest) (*models.AuthResponse, error)
	Update(email string, req *models.UpdateRequest) error
}

type authService struct {
	r repository.AuthRepository
	m *jwt.Manager
}

func NewAuthService(repo repository.AuthRepository, manager *jwt.Manager) AuthService {
	return &authService{r: repo, m: manager}
}

func (s *authService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	// --- Validation ---
	req.Email = trimmAndLowerCase(req.Email)
	if req.Email == "" || !strings.Contains(string(req.Email), "@") {
		return nil, fmt.Errorf("validation: a valid email address is required")
	}
	if len(req.Password) < 6 {
		return nil, fmt.Errorf("validation: password must be at least 6 characters long")
	}

	// --- Check for existing user ---
	existing, err := s.r.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("authService.Register: failed to check for existing admin: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("authService.Register: an admin with this email already exists.")
	}

	// --- Hash the password using bcrypt at cost 12 ---
	// Cost 12 is the modern recommendation: secure enough, not so slow it impacts UX.
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("authService.Register: failed to hash password: %w", err)
	}

	admin, err := s.r.Insert(req.Email, string(hash))
	if err != nil {
		return nil, fmt.Errorf("authService.register: failed to create admin: %w", err)
	}

	token, err := s.m.Generate(admin.ID, admin.Email)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		Admin:       admin,
	}, nil
}

func (s *authService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	req.Email = trimmAndLowerCase(req.Email)

	admin, err := s.r.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("login: database error: %w", err)
	}

	// IMPORTANT: Return the same generic error for both "user not found" and
	// "wrong password" to prevent user enumeration attacks.
	if admin == nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	token, err := s.m.Generate(admin.ID, admin.Email)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		Admin:       admin,
	}, nil
}

func (s *authService) Update(email string, req *models.UpdateRequest) error {
	admin, err := s.r.FindByEmail(email)
	if err != nil {
		return fmt.Errorf("login: database error: %w", err)
	}

	if admin == nil {
		return fmt.Errorf("auth: unknown error")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return fmt.Errorf("auth: failed to hash password: %w", err)
	}

	err = s.r.Update(email, string(hash))

	return err
}
