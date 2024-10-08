package jwtutil

type JwtUtil interface {
	ParseWithoutKey(jwt string) (map[string]any, error)
}
