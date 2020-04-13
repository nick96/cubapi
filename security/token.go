package security

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	gojwt "github.com/dgrijalva/jwt-go"
)

type JWT struct {
	claims map[string]interface{}
}

type Token struct {
	Email string
}

func (j *JWT) Subject(subject string) *JWT {
	if j.claims == nil {
		j.claims = make(map[string]interface{})
	}
	j.claims["sub"] = subject
	return j
}

func (j *JWT) Issuer(issuer string) *JWT {
	if j.claims == nil {
		j.claims = make(map[string]interface{})
	}
	j.claims["iss"] = issuer
	return j
}

func (j *JWT) Expiration(exp time.Time) *JWT {
	if j.claims == nil {
		j.claims = make(map[string]interface{})
	}
	j.claims["exp"] = strconv.FormatInt(exp.Unix(), 10)
	return j
}

func (j *JWT) ExpireIn(duration time.Duration) *JWT {
	if j.claims == nil {
		j.claims = make(map[string]interface{})
	}
	now := time.Now()
	exp := now.Add(duration)
	return j.Expiration(exp)
}

func (j *JWT) Audience(audience string) *JWT {
	if j.claims == nil {
		j.claims = make(map[string]interface{})
	}
	j.claims["aud"] = audience
	return j
}

func (j *JWT) SignedToken(secret string) (string, error) {
	token := gojwt.NewWithClaims(
		gojwt.SigningMethodHS256,
		gojwt.MapClaims(j.claims),
	)
	return token.SignedString([]byte(secret))
}

// ValidateToken validates the given token using the given secret. If the token
// is not valid in any way, an error is returned, otherwise the error is nil.
func ValidateToken(token, secret string) (Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return Token{}, fmt.Errorf("token parsing failed: %w", err)
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return Token{Email: claims["aud"].(string)}, nil
	}
	return Token{}, fmt.Errorf("token is not valid")
}
