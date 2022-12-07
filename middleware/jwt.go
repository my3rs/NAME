package middleware

import (
	"NAME/conf"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"strings"
	"sync"
	"time"
)

type UserClaims struct {
	ID       uint   `json:id`
	Username string `json:username`
	Role     string `json:"role"`
}

type JwtTokenType uint

const (
	TypeAccessToken  JwtTokenType = 1
	TypeRefreshToken JwtTokenType = 2
)

const (
	AccessTokenMaxAge  = 60 * 24 // minute // todo: time too long
	RefreshTokenMaxAge = 24      // hour
)

var (
	signer   *jwt.Signer
	verifier *jwt.Verifier
	once     sync.Once
)

func JwtMiddleware() iris.Handler {
	// init only once
	once.Do(func() {
		signer = jwt.NewSigner(jwt.HS256, conf.Config().JWT.SecretKey(), time.Minute*AccessTokenMaxAge)
		verifier = jwt.NewVerifier(jwt.HS256, conf.Config().JWT.SecretKey())
	})

	return func(ctx iris.Context) {
		VerifyAccessToken(ctx)

		ctx.Next()
	}

}

func GenerateTokenPair(userID uint, username string) (token jwt.TokenPair, err error) {
	refreshClaims := jwt.Claims{Subject: username}
	accessClaims := UserClaims{
		ID:       userID,
		Username: username,
		Role:     "admin", // todo
	}

	accessToken, err := signer.Sign(accessClaims)
	if err != nil {
		fmt.Printf("Failed to create access token: %v\n", err)
		return jwt.TokenPair{}, err
	}

	refreshToken, err := signer.Sign(refreshClaims, jwt.MaxAge(RefreshTokenMaxAge*time.Hour))
	if err != nil {
		fmt.Printf("Failed to create refresh token: %v\n", err)
		return jwt.TokenPair{}, err
	}

	return jwt.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// CheckTokenFormat check token's length and prefix
func CheckTokenFormat(token string) error {
	if len(token) == 0 {
		return errors.New("no token found in the HTTP header")
	}

	if !strings.HasPrefix(token, "Bearer ") {
		return errors.New("the format of token is not correct(\"Bearer <token>\")")
	}

	return nil
}

func VerifyAccessToken(ctx iris.Context) {
	accessToken := GetTokenFromHeader(ctx, TypeAccessToken)

	if err := CheckTokenFormat(accessToken); err != nil {
		ctx.Application().Logger().Info(err)
		ctx.StopWithJSON(401, iris.Map{"message": err.Error()})

		return
	}

	// Remove `Bearer ` prefix
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")

	verifiedToken, err := verifier.VerifyToken([]byte(accessToken))
	if err != nil {
		// Expired access token
		if err == jwt.ErrExpired {
			ctx.Application().Logger().Info("expired access token")

			ctx.StopWithJSON(440, iris.Map{"message": "expired access token"})
			return
		}

		ctx.Application().Logger().Info("invalid access token")

		ctx.StopWithJSON(401, iris.Map{"message": "invalid access token"})
		return
	}

	standardClaims := verifiedToken.StandardClaims

	_ = standardClaims
}

func GetTokenFromHeader(ctx iris.Context, t JwtTokenType) string {
	// Read token from header
	var token string
	switch t {
	case TypeAccessToken:
		token = ctx.GetHeader("Authorization")
	case TypeRefreshToken:
		token = ctx.GetHeader("Refresh-Token")
	}

	return token
}

func VerifyRefreshTokenAndGetUserName(ctx iris.Context) string {
	refreshToken := GetTokenFromHeader(ctx, TypeRefreshToken)

	if err := CheckTokenFormat(refreshToken); err != nil {
		ctx.Application().Logger().Info(err)
		ctx.StopWithJSON(401, iris.Map{"message": err.Error()})

		return ""
	}

	// Remove `Bearer ` prefix
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")

	verifiedToken, err := verifier.VerifyToken([]byte(refreshToken))
	if err != nil {
		// Expired access token
		if err == jwt.ErrExpired {
			ctx.Application().Logger().Info("expired refresh token")
			ctx.StopWithJSON(440, iris.Map{"error": "expired refresh token"})
			return ""
		}

		ctx.Application().Logger().Info("invalid refresh token")
		ctx.StopWithJSON(401, iris.Map{"message": "expired refresh token"})
		return ""
	}

	standardClaims := verifiedToken.StandardClaims

	_ = standardClaims

	return standardClaims.Subject
}
