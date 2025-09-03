package dto

import "github.com/golang-jwt/jwt/v5"

type (
	UserEvent struct {
		Event   string                 `json:"event"`
		Payload map[string]interface{} `json:"payload"`
	}

	UserCreate struct {
		Username string   `json:"user_name"`
		Password string   `json:"password"`
		Blocked  bool     `json:"blocked"`
		Scopes   []string `json:"scopes"`
	}

	UserDelete struct {
		UserName string `json:"user_name"`
	}

	UserUpdate struct {
		UserName string `json:"user_name"`
		Blocked  bool   `json:"blocked"`
	}

	UserPasswordUpdate struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}

	UserScope struct {
		UserName string `json:"user_name"`
		Scope    string `json:"scope"`
	}

	LoginRequest struct {
		UserName string `json:"user_name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	Claims struct {
		Sub     uint     `json:"sub"`
		Scopes  []string `json:"scopes"`
		Blocked bool     `json:"blocked"`
		jwt.RegisteredClaims
	}

	LoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)
