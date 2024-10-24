package helper

import (
	"hash/fnv"

	"golang.org/x/crypto/bcrypt"
)

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func CheckHash(s string, expectedHash uint32) bool {
	actualHash := Hash(s)
	return actualHash == expectedHash
}

// HashPassword generates a bcrypt hash of the password string
func HashPassword(password string) (string, error) {
	// Generate "cost" factor which determines the computational complexity
	// of the hashing. Higher the cost, more secure but slower.
	const cost = 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a password with its hashed version and returns
// true if they match, false otherwise.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
