package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// TokenClaims represents the claims stored in JWT tokens
type TokenClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	secretKey          []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

// NewJWTService creates a new JWTService instance
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey:          []byte(secretKey),
		accessTokenExpiry:  15 * time.Minute,   // Access token expires in 15 minutes
		refreshTokenExpiry: 7 * 24 * time.Hour, // Refresh token expires in 7 days
	}
}

// UserInfo contains user information for token generation
type UserInfo struct {
	UserID  string
	Email   string
	Name    string
	Picture string
}

// GenerateTokenPair generates both access and refresh tokens
func (s *JWTService) GenerateTokenPair(user UserInfo) (accessToken, refreshToken string, err error) {
	accessToken, err = s.generateToken(user, "access", s.accessTokenExpiry)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.generateToken(user, "refresh", s.refreshTokenExpiry)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// generateToken creates a JWT token with the specified claims
func (s *JWTService) generateToken(user UserInfo, tokenType string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID:    user.UserID,
		Email:     user.Email,
		Name:      user.Name,
		Picture:   user.Picture,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "go-google-auth",
			Subject:   user.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateAccessToken validates an access token
func (s *JWTService) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token
func (s *JWTService) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token from a valid refresh token
func (s *JWTService) RefreshAccessToken(refreshTokenString string) (string, error) {
	claims, err := s.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", err
	}

	user := UserInfo{
		UserID:  claims.UserID,
		Email:   claims.Email,
		Name:    claims.Name,
		Picture: claims.Picture,
	}

	return s.generateToken(user, "access", s.accessTokenExpiry)
}

// GetAccessTokenExpiry returns the access token expiry duration in seconds
func (s *JWTService) GetAccessTokenExpiry() int {
	return int(s.accessTokenExpiry.Seconds())
}

// GetRefreshTokenExpiry returns the refresh token expiry duration in seconds
func (s *JWTService) GetRefreshTokenExpiry() int {
	return int(s.refreshTokenExpiry.Seconds())
}
