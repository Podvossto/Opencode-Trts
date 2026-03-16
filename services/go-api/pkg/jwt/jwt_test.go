package jwt

import (
	"testing"
	"time"
)

const testSecret = "test-secret-key-for-unit-tests-32chars!"

func TestGenerate_Success(t *testing.T) {
	token, jti, err := Generate("user-123", "admin", testSecret, 30)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if token == "" {
		t.Fatal("Generate() returned empty token")
	}
	if jti == "" {
		t.Fatal("Generate() returned empty jti")
	}
}

func TestGenerate_DifferentJTI(t *testing.T) {
	_, jti1, _ := Generate("user-123", "admin", testSecret, 30)
	_, jti2, _ := Generate("user-123", "admin", testSecret, 30)
	if jti1 == jti2 {
		t.Fatal("Generate() should produce unique jti values")
	}
}

func TestParseRoundTrip(t *testing.T) {
	userID := "550e8400-e29b-41d4-a716-446655440000"
	role := "hr"
	expiryMinutes := 30

	token, jti, err := Generate(userID, role, testSecret, expiryMinutes)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	claims, err := Parse(token, testSecret)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("claims.UserID = %q, want %q", claims.UserID, userID)
	}
	if claims.Role != role {
		t.Errorf("claims.Role = %q, want %q", claims.Role, role)
	}
	if claims.JTI != jti {
		t.Errorf("claims.JTI = %q, want %q", claims.JTI, jti)
	}

	// ExpiresAt should be approximately 30 minutes from now
	if claims.RegisteredClaims.ExpiresAt == nil {
		t.Fatal("claims.ExpiresAt is nil")
	}
	expected := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)
	diff := claims.RegisteredClaims.ExpiresAt.Time.Sub(expected)
	if diff > 5*time.Second || diff < -5*time.Second {
		t.Errorf("ExpiresAt off by %v", diff)
	}
}

func TestParse_WrongSecret(t *testing.T) {
	token, _, err := Generate("user-123", "admin", testSecret, 30)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	_, err = Parse(token, "wrong-secret-key")
	if err == nil {
		t.Fatal("Parse() with wrong secret should return error")
	}
}

func TestParse_ExpiredToken(t *testing.T) {
	// Generate a token that expired 1 minute ago
	token, _, err := Generate("user-123", "admin", testSecret, -1)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	_, err = Parse(token, testSecret)
	if err == nil {
		t.Fatal("Parse() with expired token should return error")
	}
}

func TestParse_MalformedToken(t *testing.T) {
	_, err := Parse("not-a-real-jwt-token", testSecret)
	if err == nil {
		t.Fatal("Parse() with malformed token should return error")
	}
}

func TestParse_EmptyToken(t *testing.T) {
	_, err := Parse("", testSecret)
	if err == nil {
		t.Fatal("Parse() with empty token should return error")
	}
}
