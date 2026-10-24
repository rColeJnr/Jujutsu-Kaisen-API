package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rcolejnr/jjk-api/internal/models"
)

type jwtCustomClaims struct {
	models.Claims
	jwt.RegisteredClaims
}

// Manager handles all jwt operations using a typed struct, not a global variable.
type Manager struct {
	SecretKey   []byte
	ExpireHours int
}

func NewManager(secret string, expireHours int) *Manager {
	return &Manager{
		SecretKey:   []byte(secret),
		ExpireHours: expireHours,
	}
}

// Generate creates a signed JWT string for a given user.
func (m *Manager) Generate(id int64, email string) (string, error) {
	claims := jwtCustomClaims{
		models.Claims{ID: id,
			Email: email},
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(1) * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer:   "jjk-api"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.SecretKey)
	if err != nil {
		return "", fmt.Errorf("jwt: failed to sign token: %w", err)
	}
	return signedToken, nil
}

// Validate parses and validates a JWT string, returning the embedded claims.
func (m *Manager) Validate(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (any, error) {
		// IMPORTANT: Always verify the signing method to prevent algorithm substitution attacks
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt: unexpected signing mehtod: %v", token.Header["alg"])
		}
		return m.SecretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("jwt: token has expired")
		}
		return nil, fmt.Errorf("jwt: invalid token: %w", err)
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("jwt: token claims are invalid")
	}

	return &claims.Claims, nil
}

type contextKey string

const claimsContextKey contextKey = "jwt_claims"

// ClaimsFromContext retrieves the validated JWT Claims from the request context.
// Returns nil if not present (i.e., called on an unprotected route).
func ClaimsFromContext(ctx context.Context) *models.Claims {
	claims, _ := ctx.Value(claimsContextKey).(*models.Claims)
	return claims
}
