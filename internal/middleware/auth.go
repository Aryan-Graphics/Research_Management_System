package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"research-management-system/internal/models"
)

type contextKey string

const UserContextKey = contextKey("user")

var SecretKey = []byte("scholar-sync-research-management-secure-secret-2025")

type TokenClaims struct {
	UserID    int    `json:"uid"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"exp"`
}

// GenerateToken creates a signed token containing user claims
func GenerateToken(user *models.User) (string, error) {
	claims := TokenClaims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		ExpiresAt: time.Now().Add(48 * time.Hour).Unix(),
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadBytes)

	mac := hmac.New(sha256.New, SecretKey)
	mac.Write([]byte(payloadB64))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("%s.%s", payloadB64, signature), nil
}

// VerifyToken verifies the signature and expiration of a token
func VerifyToken(tokenStr string) (*TokenClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}

	payloadB64, sigB64 := parts[0], parts[1]

	mac := hmac.New(sha256.New, SecretKey)
	mac.Write([]byte(payloadB64))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(sigB64), []byte(expectedSig)) {
		return nil, errors.New("invalid signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, errors.New("invalid payload")
	}

	var claims TokenClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid claims")
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

// RequireAuth middleware protects routes requiring authentication
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
			return
		}

		claims, err := VerifyToken(token)
		if err != nil {
			http.Error(w, `{"error":"Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next(w, r.WithContext(ctx))
	}
}

// OptionalAuth middleware extracts user if present, but doesn't block
func OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token != "" {
			if claims, err := VerifyToken(token); err == nil {
				ctx := context.WithValue(r.Context(), UserContextKey, claims)
				next(w, r.WithContext(ctx))
				return
			}
		}
		next(w, r)
	}
}

func extractToken(r *http.Request) string {
	// 1. Authorization header: Bearer <token>
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// 2. Cookie
	if cookie, err := r.Cookie("auth_token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// 3. Query param (useful for downloads)
	if q := r.URL.Query().Get("token"); q != "" {
		return q
	}

	return ""
}

// GetUserClaims returns token claims from context if present
func GetUserClaims(ctx context.Context) *TokenClaims {
	claims, ok := ctx.Value(UserContextKey).(*TokenClaims)
	if !ok {
		return nil
	}
	return claims
}

func GetIDFromURL(r *http.Request) (int, error) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 0 {
		return 0, errors.New("missing id")
	}
	idStr := parts[len(parts)-1]
	return strconv.Atoi(idStr)
}
