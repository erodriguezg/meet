package jwtutil

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type golangJwtUtil struct{}

func NewGolangJwtUtil() JwtUtil {
	return &golangJwtUtil{}
}

// ParseWithoutKey implements JwtUtil
func (port *golangJwtUtil) ParseWithoutKey(jwtString string) (map[string]any, error) {
	token, err := jwt.Parse(jwtString, nil)
	if token == nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("can not parse claims jwt")
	}
}
