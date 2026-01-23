package helper

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

func Slugify(s string) string {
	//  Convert to lowercase
	s = strings.ToLower(s)

	//  Replace spaces and underscores with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// 3 Remove all non-alphanumeric and non-hyphen characters
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	s = reg.ReplaceAllString(s, "")

	// 4 Remove multiple consecutive hyphens
	reg2 := regexp.MustCompile(`-+`)
	s = reg2.ReplaceAllString(s, "-")

	// 5 Trim hyphens from start and end
	s = strings.Trim(s, "-")

	return s
}

func GenerateRandomApiKey(prefix string) (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("couldn't generate api key: %w", err)
	}
	randomString := hex.EncodeToString(randomBytes)
	return prefix + randomString, nil
}

func HashSHA256(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}
