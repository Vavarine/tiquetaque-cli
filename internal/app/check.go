package app

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func GetCheck(userId, name, date, time string) string {
	// 1. Inverte a string do tempo
	reversedTime := reverse(time)

	// 2. Normaliza o nome (remove espaços e coloca em minúsculas)
	normalizedName := strings.ToLower(strings.TrimSpace(name))

	// 3. Monta a string final concatenando: tempo invertido + nome + data + timestamp
	combined := reversedTime + normalizedName + date + userId

	// 4. Calcula o hash SHA-256
	hashBytes := sha256.Sum256([]byte(combined))

	// 5. Converte o hash para string hexadecimal
	return hex.EncodeToString(hashBytes[:])
}

// reverse inverte os caracteres de uma string
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
