package auth

import (
	"crypto/rand"
	"errors"
	"strings"
	"sync"
	"time"

	"NAME/conf"
	"NAME/model"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

const (
	keySize = 32 // Key size for HMAC
)

// TokenType represents different types of JWT tokens
type TokenType int

const (
	TypeAccessToken TokenType = iota
	TypeRefreshToken
)

// JWTService handles JWT token operations
type JWTService struct {
	signer    *jwt.Signer
	verifier  *jwt.Verifier
	sharedKey []byte
	config    *conf.JWTConfig
	mu        sync.RWMutex
}

var (
	instance *JWTService
	once     sync.Once
)

// GetJWTService returns a singleton instance of JWTService
func GetJWTService() *JWTService {
	once.Do(func() {
		instance = newJWTService()
	})
	return instance
}

// newJWTService creates a new JWTService instance
func newJWTService() *JWTService {
	config := conf.GetConfig().JWT
	sharedKey := make([]byte, keySize)
	_, _ = rand.Read(sharedKey)

	return &JWTService{
		signer:    jwt.NewSigner(jwt.HS256, sharedKey, time.Duration(config.AccessTokenMaxAge)),
		verifier:  jwt.NewVerifier(jwt.HS256, sharedKey),
		sharedKey: sharedKey,
		config:    &config,
	}
}

// GenerateTokenPair generates a new token pair for a user
func (s *JWTService) GenerateTokenPair(user model.User) (jwt.TokenPair, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	claims := Claims{
		ID:       user.ID,
		Username: user.Name,
		Role:     user.Role.String(),
	}

	accessToken, err := s.signer.Sign(claims)
	if err != nil {
		return jwt.TokenPair{}, err
	}

	refreshToken, err := jwt.NewSigner(jwt.HS256, s.sharedKey,
		time.Duration(s.config.RefreshTokenMaxAge)).Sign(jwt.Claims{Subject: user.Name})
	if err != nil {
		return jwt.TokenPair{}, err
	}

	return jwt.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// VerifyToken verifies a token and returns the claims
func (s *JWTService) VerifyToken(tokenString string) (*Claims, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.checkTokenFormat(tokenString); err != nil {
		return nil, err
	}

	// Remove Bearer prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	verifiedToken, err := s.verifier.VerifyToken([]byte(tokenString))
	if err != nil {
		if err == jwt.ErrExpired {
			return nil, errors.New("token expired")
		}
		return nil, errors.New("invalid token")
	}

	claims := new(Claims)
	if err := verifiedToken.Claims(&claims); err != nil {
		return nil, err
	}

	return claims, nil
}

// checkTokenFormat validates the token format
func (s *JWTService) checkTokenFormat(token string) error {
	if token == "" {
		return errors.New("missing token")
	}
	if !strings.HasPrefix(token, "Bearer ") {
		return errors.New("invalid token format")
	}
	return nil
}

// GetTokenFromHeader extracts token from request header
func (s *JWTService) GetTokenFromHeader(ctx iris.Context, t TokenType) string {
	var authHeader string
	switch t {
	case TypeAccessToken:
		authHeader = ctx.GetHeader("Authorization")
	case TypeRefreshToken:
		authHeader = ctx.GetHeader("X-Refresh-Token")
	}
	return authHeader
}

// GetClaimsFromContext extracts claims from context
func (s *JWTService) GetClaimsFromContext(ctx iris.Context) *Claims {
	if claims := ctx.Values().Get("jwt").(*jwt.VerifiedToken); claims != nil {
		userClaims := new(Claims)
		if err := claims.Claims(userClaims); err != nil {
			return nil
		}
		return userClaims
	}
	return nil
}
