package auth

import "github.com/kataras/iris/v12"

// JWTMiddleware creates a new JWT middleware handler
func JWTMiddleware() iris.Handler {
	jwtService := GetJWTService()

	return func(ctx iris.Context) {
		token := jwtService.GetTokenFromHeader(ctx, TypeAccessToken)
		claims, err := jwtService.VerifyToken(token)
		if err != nil {
			ctx.StopWithJSON(401, iris.Map{
				"message": err.Error(),
			})
			return
		}

		// Store claims in context for later use
		ctx.Values().Set("claims", claims)
		ctx.Next()
	}
}
