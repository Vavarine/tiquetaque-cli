package app

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// GetCheck generates a SHA-256 hash based on the provided parameters. This is a reverse-engineered
// version of the check generation algorithm used by the TiqueTaque web app.
// It combines the reversed time, normalized name, date, and userId to create a unique hash.
func GetCheck(userId, name, date, time string) string {
	reversedTime := reverse(time)
	normalizedName := strings.ToLower(strings.TrimSpace(name))
	combined := reversedTime + normalizedName + date + userId
	hashBytes := sha256.Sum256([]byte(combined))

	return hex.EncodeToString(hashBytes[:])
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
