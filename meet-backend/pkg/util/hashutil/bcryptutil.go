package hashutil

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"encoding/base64"
)

func BCryptHash(seed string) (string, error) {

	sha256Hash := sha256.Sum256([]byte(seed))

	hashBytes, err := bcrypt.GenerateFromPassword(sha256Hash[:], bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("error at bcrypt.GenerateFromPassword. error: %w", err)
	}
	return string(hashBytes), nil
}

func B64UrlEncoding(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}

func SHA256HashB64UrlEncoding(data string) string {
	sha256Hash := sha256.Sum256([]byte(data))
	return base64.URLEncoding.EncodeToString(sha256Hash[:])
}
