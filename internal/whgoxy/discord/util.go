package discord

import "math/rand"

//goland:noinspection SpellCheckingInspection
const SecretAvailableChars string = "ABCDEFGHKMNPQRSTUVXYZabcdefghmnpqrstuvwxyz023456789._-"

// GenerateSecret generates a 64 long string
func GenerateSecret() (res string) {
	var secret [64]byte

	l := len(SecretAvailableChars)
	for i := 0; i < 64; i++ {
		idx := rand.Intn(l)
		secret[i] = SecretAvailableChars[idx]
	}

	return string(secret[:])
}
