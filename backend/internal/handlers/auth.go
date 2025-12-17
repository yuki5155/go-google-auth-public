package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuki5155/go-google-auth/internal/config"
	"github.com/yuki5155/go-google-auth/internal/services"
	"google.golang.org/api/idtoken"
)

// GoogleAuthRequest represents the request body for Google authentication
type GoogleAuthRequest struct {
	Credential string `json:"credential" binding:"required"`
}

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	config     *config.Config
	jwtService *services.JWTService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(cfg *config.Config, jwtService *services.JWTService) *AuthHandler {
	return &AuthHandler{
		config:     cfg,
		jwtService: jwtService,
	}
}

// GoogleLogin handles Google Sign-In token verification and JWT generation
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Missing or invalid credential",
		})
		return
	}

	// Verify the Google ID token
	ctx := context.Background()
	payload, err := idtoken.Validate(ctx, req.Credential, h.config.GoogleClientID)
	if err != nil {
		log.Printf("Failed to verify Google ID token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_token",
			"message": "Failed to verify Google ID token",
		})
		return
	}

	// Extract user information from the token payload
	userID, _ := payload.Claims["sub"].(string)
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)

	if !emailVerified {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unverified_email",
			"message": "Email address is not verified",
		})
		return
	}

	// Generate JWT tokens
	userInfo := services.UserInfo{
		UserID:  userID,
		Email:   email,
		Name:    name,
		Picture: picture,
	}

	accessToken, refreshToken, err := h.jwtService.GenerateTokenPair(userInfo)
	if err != nil {
		log.Printf("Failed to generate JWT tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "token_generation_failed",
			"message": "Failed to generate authentication tokens",
		})
		return
	}

	// Set cookies
	secure := h.config.IsProduction()
	sameSite := http.SameSiteLaxMode

	// Access token cookie
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   h.jwtService.GetAccessTokenExpiry(),
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	}
	http.SetCookie(c.Writer, accessCookie)

	// Refresh token cookie
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   h.jwtService.GetRefreshTokenExpiry(),
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	}
	http.SetCookie(c.Writer, refreshCookie)

	log.Printf("User logged in successfully: %s (%s)", email, userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":      userID,
			"email":   email,
			"name":    name,
			"picture": picture,
		},
	})
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "missing_refresh_token",
			"message": "Refresh token not found",
		})
		return
	}

	// Generate new access token from refresh token
	newAccessToken, err := h.jwtService.RefreshAccessToken(refreshToken)
	if err != nil {
		if err == services.ErrExpiredToken {
			// Clear cookies if refresh token is expired
			h.clearAuthCookies(c)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "refresh_token_expired",
				"message": "Refresh token has expired, please login again",
			})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_refresh_token",
			"message": "Invalid refresh token",
		})
		return
	}

	// Set new access token cookie
	secure := h.config.IsProduction()
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		Path:     "/",
		MaxAge:   h.jwtService.GetAccessTokenExpiry(),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(c.Writer, accessCookie)

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	h.clearAuthCookies(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// GetCurrentUser returns the current authenticated user's information
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// User information is set by the auth middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	email, _ := c.Get("email")
	name, _ := c.Get("name")
	picture, _ := c.Get("picture")

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":      userID,
			"email":   email,
			"name":    name,
			"picture": picture,
		},
	})
}

// clearAuthCookies removes authentication cookies
func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	secure := h.config.IsProduction()

	// Clear access token
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
	})

	// Clear refresh token
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
	})
}
