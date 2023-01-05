package middleware

import (
	"NAME/conf"
	"NAME/model"
	"NAME/service"
	"crypto/rsa"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"strconv"
	"sync"
	"time"
)

type UserClaims struct {
	ID       uint           `json:"id"`
	Username string         `json:"username"`
	Role     model.UserRole `json:"role"`
}

func (u *UserClaims) GetID() uint {
	return u.ID
}

func (u *UserClaims) GetUsername() string {
	return u.Username
}

type JwtTokenType uint

const (
	TypeAccessToken  JwtTokenType = 1
	TypeRefreshToken JwtTokenType = 2
)

const (
	accessTokenMaxAge  = 60 * 24 * time.Minute
	refreshTokenMaxAge = 7 * 24 * 60 * time.Minute // hour
)

//var (
//	signer   *jwt.Signer
//	verifier *jwt.Verifier
//	once     sync.Once
//)

var (
	once sync.Once
	//privateKey, publicKey = jwt.MustLoadRSA(conf.Config().JWT.PrivateKey, conf.Config().JWT.PublicKey)
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	signer     *jwt.Signer
	verifier   *jwt.Verifier
)

//func JwtMiddleware() iris.Handler {
//	// init only once
//	once.Do(func() {
//		signer = jwt.NewSigner(jwt.HS256, conf.Config().JWT.SecretKey(), time.Minute*AccessTokenMaxAge)
//		verifier = jwt.NewVerifier(jwt.HS256, conf.Config().JWT.SecretKey())
//	})
//
//	return func(ctx iris.Context) {
//		VerifyAccessToken(ctx)
//
//		ctx.Next()
//	}
//}

func JwtMiddleware() iris.Handler {
	once.Do(func() {
		privateKey, publicKey = jwt.MustLoadRSA(conf.Config().JWT.PrivateKey, conf.Config().JWT.PublicKey)
		signer = jwt.NewSigner(jwt.RS256, privateKey, accessTokenMaxAge)
		verifier = jwt.NewVerifier(jwt.RS256, publicKey)
	})

	return verifier.Verify(func() interface{} {
		return new(UserClaims)
	})
}

func GenerateTokenPair(user model.User) (token jwt.TokenPair, err error) {
	JwtMiddleware()
	// Map the current user with the refresh token,
	// so we make sure, on refresh route, that this refresh token owns
	// to that user before re-generate.
	refreshClaims := jwt.Claims{Subject: user.Name}
	accessClaims := UserClaims{
		ID:       user.ID,
		Username: user.Name,
		Role:     user.Role,
	}

	return signer.NewTokenPair(accessClaims, refreshClaims, refreshTokenMaxAge)
}

func RefreshToken(ctx iris.Context) {
	JwtMiddleware()

	refreshToken := getRefreshToken(ctx)

	verifiedToken, err := verifier.VerifyToken([]byte(refreshToken), jwt.Expected{})
	if err != nil {
		ctx.StopWithError(iris.StatusNonAuthoritativeInfo, err)
		return
	}

	userID, err := strconv.Atoi(verifiedToken.StandardClaims.Subject)
	if err != nil {
		ctx.StopWithError(iris.StatusNonAuthoritativeInfo, err)
		return
	}
	user, err := service.GetUserService().GetUserById(userID)
	if err != nil {
		ctx.StopWithError(iris.StatusNonAuthoritativeInfo, err)
		return
	}

	pair, err := GenerateTokenPair(user)
	if err != nil {
		ctx.StopWithError(iris.StatusNonAuthoritativeInfo, err)
		return
	}

	ctx.StatusCode(200)
	ctx.Header("authorization", string(pair.AccessToken))
	ctx.Header("refresh-token", string(pair.RefreshToken))
	ctx.JSON(iris.Map{"message": "refresh success: check new tokens in the HTTP header"})

	ctx.Next()
}

func getRefreshToken(ctx iris.Context) string {
	return ctx.GetHeader("Refresh-Token")
}

func getAccessToken(ctx iris.Context) string {
	return ctx.GetHeader("Authorization")
}

// CheckTokenFormat check token's length and prefix
//func CheckTokenFormat(token string) error {
//	if len(token) == 0 {
//		return errors.New("no token found in the HTTP header")
//	}
//
//	if !strings.HasPrefix(token, "Bearer ") {
//		return errors.New("the format of token is not correct(\"Bearer <token>\")")
//	}
//
//	return nil
//}

//func VerifyAccessToken(ctx iris.Context) {
//	accessToken := GetTokenFromHeader(ctx, TypeAccessToken)
//
//	if err := CheckTokenFormat(accessToken); err != nil {
//		ctx.Application().Logger().Info(err)
//		ctx.StopWithJSON(401, iris.Map{"message": err.Error()})
//
//		return
//	}
//
//	// Remove `Bearer ` prefix
//	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
//
//	verifiedToken, err := verifier.VerifyToken([]byte(accessToken))
//	if err != nil {
//		// Expired access token
//		if err == jwt.ErrExpired {
//			ctx.Application().Logger().Info("expired access token")
//
//			ctx.StopWithJSON(440, iris.Map{"message": "expired access token"})
//			return
//		}
//
//		ctx.Application().Logger().Info("invalid access token")
//
//		ctx.StopWithJSON(401, iris.Map{"message": "invalid access token"})
//		return
//	}
//
//	standardClaims := verifiedToken.StandardClaims
//
//	_ = standardClaims
//}
//
//func GetTokenFromHeader(ctx iris.Context, t JwtTokenType) string {
//	// Read token from header
//	var token string
//	switch t {
//	case TypeAccessToken:
//		token = ctx.GetHeader("Authorization")
//	case TypeRefreshToken:
//		token = ctx.GetHeader("Refresh-Token")
//	}
//
//	return token
//}
//
//func VerifyRefreshTokenAndGetUserName(ctx iris.Context) string {
//	refreshToken := GetTokenFromHeader(ctx, TypeRefreshToken)
//
//	if err := CheckTokenFormat(refreshToken); err != nil {
//		ctx.Application().Logger().Info(err)
//		ctx.StopWithJSON(401, iris.Map{"message": err.Error()})
//
//		return ""
//	}
//
//	// Remove `Bearer ` prefix
//	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")
//
//	verifiedToken, err := verifier.VerifyToken([]byte(refreshToken))
//	if err != nil {
//		// Expired access token
//		if err == jwt.ErrExpired {
//			ctx.Application().Logger().Info("expired refresh token")
//			ctx.StopWithJSON(440, iris.Map{"error": "expired refresh token"})
//			return ""
//		}
//
//		ctx.Application().Logger().Info("invalid refresh token")
//		ctx.StopWithJSON(401, iris.Map{"message": "expired refresh token"})
//		return ""
//	}
//
//	standardClaims := verifiedToken.StandardClaims
//
//	_ = standardClaims
//
//	return standardClaims.Subject
//}
