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
	private, public := jwt.MustLoadRSA(conf.GetConfig().JWT.PrivateKey, conf.GetConfig().JWT.PublicKey)

	return &JWTService{
		signer:        jwt.NewSigner(jwt.RS256, private, config.AccessTokenMaxAge),
		accessSigner:  jwt.NewSigner(jwt.RS256, private, config.AccessTokenMaxAge),
		refreshSigner: jwt.NewSigner(jwt.RS256, private, config.RefreshTokenMaxAge),
		verifier:      jwt.NewVerifier(jwt.RS256, public),
		config:        &config,
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

	fmt.Printf("access token: %s\n", tokenPair.AccessToken)

	return tokenPair, nil
}

func (s *JWTService) VerifyAccessToken(ctx iris.Context) error {
	accessToken := s.GetTokenFromHeader(ctx, dict.TypeAccessToken)

	if accessToken == "" {
		return customerror.NewJWTInvalidError("empty access token")
	}

	if err := s.checkTokenFormat(accessToken); err != nil {
		return customerror.NewJWTFormatError("invalid token format")
	}

	if err := s.checkTokenFormat(accessToken); err != nil {
		return customerror.NewJWTFormatError("invalid token format")
	}

	// 去掉Bearer前缀
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	fmt.Println(accessToken)

	token, err := s.verifier.VerifyToken([]byte(accessToken))
	if err != nil {
		fmt.Println(token)
		return customerror.NewJWTInvalidError("invalid token: " + string(err.Error()))
	}

	if conf.GetConfig().Mode == conf.DEV {
		log.Println("[DEBUG] token verified: ", token)
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
	accessToken := s.GetTokenFromHeader(ctx, dict.TypeAccessToken)
	token, err := s.verifier.VerifyToken([]byte(accessToken))
	if err != nil {
		return model.Claims{}, customerror.NewJWTInvalidError("invalid token")
	}

	// Get the custom claims from the token
	var claims model.Claims
	if err := token.Claims(&claims); err != nil {
		return model.Claims{}, customerror.NewJWTInvalidError("invalid token: " + string(err.Error()))
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
