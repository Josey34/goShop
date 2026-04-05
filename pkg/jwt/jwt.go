package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	CustomerID string
	Email      string
	jwt.RegisteredClaims
}

type JWTService struct {
	secret string
	expiry int
}

func NewJWTService(secret string, expiry int) *JWTService {
	return &JWTService{
		secret: secret,
		expiry: expiry,
	}
}

func (j *JWTService) Generate(customerID, email string) (string, error) {
	claims := &Claims{
		CustomerID: customerID,
		Email:      email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.expiry))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secret))
}

func (j *JWTService) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token claims")
}
