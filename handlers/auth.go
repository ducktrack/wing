package handlers

import (
	"bytes"
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"strings"
	"time"
)

// CreateToken creates a JWT token
func CreateToken(recordToken RecordToken, duration time.Duration) (*jwt.Token, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(recordToken)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to encode RecordToken to JSON")
	}

	jsonPayload := string(buffer.Bytes())

	return jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Issuer:    "wing",
		ExpiresAt: time.Now().Add(duration).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   jsonPayload,
	}), nil
}

// SignToken signs the token
func SignToken(token *jwt.Token, secret string) (string, error) {
	return token.SignedString([]byte(secret))
}

// ParseAndVerifyAuthenticationHeader definition
func ParseAndVerifyAuthenticationHeader(headerToken string, secret string) (*jwt.Token, error) {
	tokenString, err := stripBearerPrefixFromTokenString(headerToken)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to strip 'Beared' prefix from token ('%s')", headerToken)
	}

	return ParseAndVerifyToken(tokenString, secret)
}

// ParseAndVerifyToken parses and verifies token
func ParseAndVerifyToken(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// validate if the alg is what is expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if token.Valid {
		return token, nil
	}

	if validationError, ok := err.(*jwt.ValidationError); ok {
		if validationError.Errors&(jwt.ValidationErrorMalformed) != 0 {
			return nil, errors.Wrapf(err, "Malformed token ('%s')", tokenString)

		} else if validationError.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, errors.Wrapf(err, "Token expired or inactive")
		}

		return nil, errors.Wrapf(err, "Failed to parse token ('%s')", tokenString)
	}

	return nil, errors.Errorf("Failed to parse token ('%s')", tokenString)
}

// Strips 'Bearer ' prefix from bearer token string
func stripBearerPrefixFromTokenString(token string) (string, error) {
	// Should be a bearer token
	if len(token) > 6 && strings.ToUpper(token[0:7]) == "BEARER " {
		return token[7:], nil
	}

	return "", errors.Errorf("Token is not a 'Bearer' token")
}
