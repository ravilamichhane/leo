package jwt

import (
	"time"

	"github.com/ravilmc/leo/types"
	"github.com/ravilmc/leo/web"

	gojwt "github.com/golang-jwt/jwt/v5"
)

func CreateJWT(id uint, email string, role types.UserRole) (*types.Tokens, error) {
	expirationTime := web.GetEnvInt("JWT_ACCESS_TOKEN_TTL", 900)
	refreshTokenExpirationTime := web.GetEnvInt("JWT_REFRESH_TOKEN_TTL", 86400)
	claims := &types.CustomClaims{
		ID:    id,
		Email: email,
		Role:  role,
		RegisteredClaims: gojwt.RegisteredClaims{
			IssuedAt:  gojwt.NewNumericDate(time.Now()),
			Audience:  gojwt.ClaimStrings{web.GetEnv("JWT_TOKEN_AUDIENCE", "service")},
			Issuer:    web.GetEnv("JWT_TOKEN_ISSUER", "service"),
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Duration(expirationTime) * time.Second)),
		},
	}

	refreshTokenClaims := &types.RefreshTokenClaims{
		ID: id,
		RegisteredClaims: gojwt.RegisteredClaims{
			IssuedAt:  gojwt.NewNumericDate(time.Now()),
			Audience:  gojwt.ClaimStrings{web.GetEnv("JWT_TOKEN_AUDIENCE", "service")},
			Issuer:    web.GetEnv("JWT_TOKEN_ISSUER", "service"),
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Duration(refreshTokenExpirationTime) * time.Second)),
		},
	}

	refreshToken := gojwt.NewWithClaims(gojwt.SigningMethodHS256, refreshTokenClaims)

	refreshTokenSecret := web.GetEnv("JWT_REFRESH_TOKEN_SECRET", "refresh_secret")

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)

	secret := web.GetEnv("JWT_SECRET", "secret")

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	refreshTokenString, err := refreshToken.SignedString([]byte(refreshTokenSecret))
	if err != nil {
		return nil, err
	}

	return &types.Tokens{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func VerifyJWT(tokenString string) (*types.CustomClaims, error) {
	secret := web.GetEnv("JWT_SECRET", "secret")

	var customClaims types.CustomClaims

	token, err := gojwt.ParseWithClaims(tokenString, &customClaims, func(token *gojwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {

		return nil, err
	}

	claims, ok := token.Claims.(*types.CustomClaims)
	if !ok {

		return nil, err
	}

	return claims, nil
}

func VerifyRefreshToken(tokenString string) (*types.RefreshTokenClaims, error) {
	secret := web.GetEnv("JWT_REFRESH_TOKEN_SECRET", "refresh_secret")

	var customClaims types.RefreshTokenClaims

	_, err := gojwt.ParseWithClaims(tokenString, &customClaims, func(token *gojwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	return &customClaims, nil
}
