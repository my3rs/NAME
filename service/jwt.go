package service

import (
	"NAME/customerror"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"NAME/conf"
	"NAME/dict"
	"NAME/model"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

// JWTService handles JWT token operations
type JWTService struct {
	accessSigner  *jwt.Signer
	refreshSigner *jwt.Signer
	signer        *jwt.Signer
	verifier      *jwt.Verifier
	config        *conf.JWTConfig
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

	return &JWTService{
		signer:   jwt.NewSigner(jwt.HS256, config.Secret, config.AccessTokenMaxAge),
		verifier: jwt.NewVerifier(jwt.HS256, config.Secret),
		config:   &config,
	}
}

// GenerateTokenPair generates a new token pair for a user
func (s *JWTService) GenerateTokenPair(user model.User) (jwt.TokenPair, error) {
	now := time.Now()

	// Create refresh claims with user ID as subject
	refreshClaims := jwt.Claims{
		Subject:  fmt.Sprintf("%s", user.Username),
		Issuer:   "NAME",
		IssuedAt: now.Unix(),
		Expiry:   now.Add(time.Second * s.config.RefreshTokenMaxAge).Unix(),
	}

	// Create access claims with user details
	accessClaims := model.Claims{
		Claims: jwt.Claims{
			Subject:  fmt.Sprintf("%s", user.Username),
			Issuer:   "NAME",
			IssuedAt: now.Unix(),
			Expiry:   now.Add(time.Second * s.config.AccessTokenMaxAge).Unix(),
		},
		Role: user.Role,
	}

	tokenPair, err := s.signer.NewTokenPair(accessClaims, refreshClaims, s.config.RefreshTokenMaxAge)
	if err != nil {
		return jwt.TokenPair{}, err
	}

	return tokenPair, nil
}

func (s *JWTService) VerifyAccessToken(ctx iris.Context) error {
	accessToken := s.GetTokenFromHeader(ctx, dict.TypeAccessToken)

	if accessToken == "" {
		return fmt.Errorf("access token is empty")
	}

	if err := s.checkTokenFormat(accessToken); err != nil {
		return fmt.Errorf("invalid access token: %w", err)
	}

	// 去掉Bearer前缀
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")

	_, err := s.verifier.VerifyToken([]byte(accessToken))
	if err != nil {
		log.Println("invalid access token: %w", err)
		return fmt.Errorf("invalid access token: %w", err)
	}

	return nil
}

func (s *JWTService) VerifyRefreshToken(ctx iris.Context) (model.Claims, error) {
	refreshToken := s.GetTokenFromHeader(ctx, dict.TypeRefreshToken)

	token, err := s.verifier.VerifyToken([]byte(refreshToken))
	if err != nil {
		return model.Claims{}, customerror.NewJWTInvalidError("invalid token: " + string(err.Error()))
	}

	if conf.GetConfig().Mode == conf.DEV {
		log.Println("[DEBUG] token verified: ", token)
	}

	// Get the custom claims from the token
	var claims model.Claims
	if err := token.Claims(&claims); err != nil {
		return model.Claims{}, customerror.NewJWTInvalidError("invalid token: " + string(err.Error()))
	}

	return claims, nil
}

func (s *JWTService) GetClaimsFromContext(ctx iris.Context) (model.Claims, error) {
	// 从请求头获取 token
	token := s.GetTokenFromHeader(ctx, dict.TypeAccessToken)
	if err := s.checkTokenFormat(token); err != nil {
		return model.Claims{}, customerror.NewJWTInvalidError(err.Error())
	}

	// 去掉 "Bearer " 前缀
	tokenString := token[7:]

	// 验证 token
	verifiedToken, err := s.verifier.VerifyToken([]byte(tokenString))
	if err != nil {
		return model.Claims{}, customerror.NewJWTInvalidError("invalid token: " + err.Error())
	}

	// 解析 claims
	var claims model.Claims
	if err := verifiedToken.Claims(&claims); err != nil {
		return model.Claims{}, customerror.NewJWTInvalidError("invalid claims: " + err.Error())
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
func (s *JWTService) GetTokenFromHeader(ctx iris.Context, t dict.TokenType) string {
	var authHeader string
	switch t {
	case dict.TypeAccessToken:
		authHeader = ctx.GetHeader("Authorization")
	case dict.TypeRefreshToken:
		authHeader = ctx.GetHeader("X-Refresh-Token")
	}
	return authHeader
}
