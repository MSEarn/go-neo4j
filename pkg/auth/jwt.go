package auth

import (
	"context"
	"time"

	"github.com/MSEarn/go-neo4j/pkg/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type JWT struct {
	SigningMethod jwt.SigningMethod
	SigningKey    interface{}
	TTL           time.Duration
}

func NewJWT(method string, secret interface{}, ttl string) (*JWT, error) {
	duration, err := time.ParseDuration(ttl)
	if err != nil {
		return nil, err
	}

	return &JWT{
		SigningMethod: jwt.GetSigningMethod(method),
		SigningKey:    secret,
		TTL:           duration,
	}, nil
}

type SignFunc func(ctx context.Context, subject string, data map[string]interface{}) (string, error)

func NewSignFunc(j *JWT) SignFunc {
	return func(ctx context.Context, subject string, data map[string]interface{}) (string, error) {
		claims := j.createClaims(subject, data)
		token := jwt.NewWithClaims(j.SigningMethod, claims)
		return token.SignedString(j.SigningKey)
	}
}

func (j *JWT) createClaims(subject string, data map[string]interface{}) *jwt.MapClaims {
	claims := jwt.MapClaims{}
	for key, value := range data {
		claims[key] = value
	}
	claims["sub"] = subject
	claims["exp"] = util.Now().Add(j.TTL).Unix()

	return &claims
}

type VerifyFunc func(ctx context.Context, tokenString string) (*jwt.Token, *jwt.MapClaims, error)

func NewVerify(j *JWT) VerifyFunc {
	return func(ctx context.Context, tokenString string) (*jwt.Token, *jwt.MapClaims, error) {
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method != j.SigningMethod {
				return nil, errors.Errorf("Unexpected jwt signing method = %v", token.Header["alg"])
			}
			return j.SigningKey, nil
		})
		if err != nil || !token.Valid {
			return nil, nil, err
		}
		if err = claims.Valid(); err != nil {
			return nil, nil, err
		}
		return token, &claims, nil
	}
}
