package discord

import "math/rand"

const SecretAvailableChars string = "ABCDEFGHKMNPQRSTUVXYZabcdefghmnpqrstuvwxyz023456789."

// generateSecret generates a 64 long string
func generateSecret() (res string) {
	var secret [64]byte

	l := len(SecretAvailableChars)
	for i := 0; i < 64; i++ {
		idx := rand.Intn(l)
		secret[i] = SecretAvailableChars[idx]
	}

	return string(secret[:])
}
