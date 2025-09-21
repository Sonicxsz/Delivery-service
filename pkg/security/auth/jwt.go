package security

import (
	"context"
	"errors"
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	SecretJWTKey string `toml:"jwt_secret_key"`
	Audience     string `toml:"jwt_audience"`
	Issuer       string `toml:"jwt_issuer"`
}

func NewJWTConfig() *JWTConfig {
	return &JWTConfig{}
}

func (j *JWTConfig) emptyFunc(context.Context) (any, error) {
	return []byte(j.SecretJWTKey), nil
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (c *CustomClaims) Validate(ctx context.Context) error {

	if c.UserID == "" {
		return errors.New("user_id cannot be empty")
	}
	return nil
}

func NewJwtMiddleware(config *JWTConfig) *jwtmiddleware.JWTMiddleware {
	var jwtValidator, err = validator.New(
		config.emptyFunc,
		validator.HS256,
		config.Issuer,
		[]string{config.Audience},
		validator.WithCustomClaims(func() validator.CustomClaims {
			return &CustomClaims{}
		}),
	)

	if err != nil {
		println("Something went wrong while configuring JWT middleware", err.Error())
	}

	return jwtmiddleware.New(jwtValidator.ValidateToken)
}

func GetClaimsFromContext(r *http.Request) (*CustomClaims, error) {
	token := r.Context().Value(jwtmiddleware.ContextKey{})

	if token == nil {
		return nil, errors.New("token not found")
	}

	claims, ok := token.(*validator.ValidatedClaims)
	if !ok {
		return nil, errors.New("invalid claims type (expected ValidatedClaims)")
	}

	customClaims, ok := claims.CustomClaims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid custom claims type")
	}

	return customClaims, nil
}

func GenerateJWT(id string, jwtConfig *JWTConfig) (string, error) {
	claims := CustomClaims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			Issuer:    jwtConfig.Issuer,
			Audience:  []string{jwtConfig.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtConfig.SecretJWTKey))
}
