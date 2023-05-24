package utils

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var key = os.Getenv("JWT_SECRET")

const tokenLifespan = 15 * 60 // 15 minutes
const leeway = 3 * 60         // 5 minutes
const issuer = "travelUtilities-expenseApi"

func GenerateToken(userId *uuid.UUID) (string, error) {
	now := time.Now()

	claims := &jwt.MapClaims{
		"exp": now.Add(time.Duration(tokenLifespan) * time.Second).Unix(),
		"iat": now.Unix(),
		"iss": issuer,
		"sub": userId.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func ValidateToken(tokenString string) (*uuid.UUID, error) {
	validMethodsOption := jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})
	leewayOption := jwt.WithLeeway(time.Duration(leeway) * time.Second)
	issuerOption := jwt.WithIssuer(issuer)
	issuedAtOption := jwt.WithIssuedAt()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}, validMethodsOption, leewayOption, issuerOption, issuedAtOption)

	if err != nil {
		log.Printf("Error in jwtutils.ValidateToken().jwt.Parse(): %v", err.Error())
		return nil, err
	}

	var claims jwt.MapClaims
	var ok bool

	if claims, ok = token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		log.Printf("Error in jwtutils.ValidateToken().token.Claims(): %v", err.Error())
		return nil, err
	}

	userId, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		log.Printf("Error in jwtutils.ValidateToken().uuid.Parse(): %v", err.Error())
		return nil, err
	}

	return &userId, nil
}

func ExtractToken(authHeader string) (string, error) {
	if len(authHeader) < 7 {
		return "", jwt.ErrInvalidKey
	}

	// Return token without "Bearer " prefix
	return authHeader[7:], nil
}
