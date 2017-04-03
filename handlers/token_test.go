package handlers_test

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/duckclick/wing/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"time"
)

type TokenTestSuite struct {
	suite.Suite
	payload     string
	secret      string
	signedToken string
	duration    time.Duration
}

func createSignedToken(suite *TokenTestSuite) string {
	token := handlers.CreateToken(suite.payload, suite.duration)
	tokenString, _ := handlers.SignToken(token, suite.secret)
	return tokenString
}

func (suite *TokenTestSuite) SetupTest() {
	suite.payload = `{"secret": "value"}`
	suite.secret = "my-secret"
	suite.duration = time.Duration(10) * time.Second
}

func (suite *TokenTestSuite) TestCreateToken() {
	token := handlers.CreateToken(suite.payload, suite.duration)
	assert.NotNil(suite.T(), token)

	claims := token.Claims.(*jwt.StandardClaims)
	assert.NotNil(suite.T(), claims)
	assert.Equal(suite.T(), claims.Issuer, "wing")
	assert.Equal(suite.T(), claims.ExpiresAt, time.Now().Add(suite.duration).Unix())
	assert.Equal(suite.T(), claims.IssuedAt, time.Now().Unix())
	assert.Equal(suite.T(), claims.Subject, suite.payload)
}

func (suite *TokenTestSuite) TestSignToken() {
	token := handlers.CreateToken(suite.payload, suite.duration)
	assert.NotNil(suite.T(), token)

	tokenString, err := handlers.SignToken(token, suite.secret)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), tokenString)

	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(_ *jwt.Token) (interface{}, error) {
		return []byte(suite.secret), nil
	})
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), parsedToken)
	claims := token.Claims.(*jwt.StandardClaims)
	assert.NotNil(suite.T(), claims)
	assert.Equal(suite.T(), claims.Subject, suite.payload)
}

func (suite *TokenTestSuite) TestParseAndVerifyToken() {
	signedToken := createSignedToken(suite)
	token, err := handlers.ParseAndVerifyToken(signedToken, suite.secret)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), token)
}

func (suite *TokenTestSuite) TestParseAndVerifyTokenWhenTokenIsExpired() {
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTEwNTA2NjcsImlhdCI6MTQ5MTA1MDY1NywiaXNzIjoid2luZyIsInN1YiI6IntcInNlY3JldFwiOiBcInZhbHVlXCJ9In0.KyVR9aw33DUB0ybpqY3XuDxWlbLXGSk1CixTNk0SBDQ"
	token, err := handlers.ParseAndVerifyToken(expiredToken, suite.secret)
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(err.Error(), "Token expired or inactive"), "Token should be expired")
	assert.Nil(suite.T(), token)
}

func (suite *TokenTestSuite) TestParseAndVerifyTokenWhenTokenHasAInvalidSignature() {
	wrongSecretToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ"
	token, err := handlers.ParseAndVerifyToken(wrongSecretToken, suite.secret)
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(err.Error(), "Failed to parse token"), "Token should have invalid signature")
	assert.True(suite.T(), strings.Contains(err.Error(), "signature is invalid"), "Token should have invalid signature")
	assert.Nil(suite.T(), token)
}

func (suite *TokenTestSuite) TestParseAndVerifyTokenWhenTokenIsMalformed() {
	malformedToken := "eyJhbGciJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiITY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFON"
	token, err := handlers.ParseAndVerifyToken(malformedToken, suite.secret)
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(err.Error(), "Malformed token"), "Token should be malformed")
	assert.Nil(suite.T(), token)
}

func (suite *TokenTestSuite) TestParseAndVerifyTokenWhenTokenWasSignedWithADifferentAlgorithm() {
	tokenWithDifferentAlgorithm := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.EkN-DOsnsuRjRO6BxXemmJDm3HbxrbRzXglbN2S4sOkopdU4IsDxTI8jO19W_A4K8ZPJijNLis4EZsHeY559a4DFOd50_OqgHGuERTqYZyuhtF39yxJPAjUESwxk2J5k_4zM3O-vtd1Ghyo4IbqKKSy6J9mTniYJPenn5-HIirE"
	token, err := handlers.ParseAndVerifyToken(tokenWithDifferentAlgorithm, suite.secret)
	assert.NotNil(suite.T(), err)
	assert.True(suite.T(), strings.Contains(err.Error(), "Unexpected signing method"), "Token should have been signed with a different algorithm")
	assert.Nil(suite.T(), token)
}

func TestToken(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}
