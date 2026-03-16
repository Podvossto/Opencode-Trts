// Purpose: JWT helper — token generation, parsing, claims extraction
// Ref: ADR-002 — JWT with jti claim for Redis blacklist
// Dev must implement: Generate(userID, role, jti string) string, Parse(token string) (*Claims, error)
package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims is the JWT payload structure
type Claims struct {
	UserID string `json:"sub"`
	Role   string `json:"role"`
	JTI    string `json:"jti"` // Used for Redis blacklist lookup
	jwt.RegisteredClaims
}

// Generate creates a signed JWT for the given user
func Generate(userID, role, secret string, expiryMinutes int) (string, string, error) {
	jti := uuid.New().String()
	claims := Claims{
		UserID: userID,
		Role:   role,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", "", fmt.Errorf("jwt sign error: %w", err)
	}

	return signed, jti, nil
}

// Parse validates and parses a JWT string, returning Claims
func Parse(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt parse error: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
