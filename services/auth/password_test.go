package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEditUser(t *testing.T) {
	t.Run("Hash password, should return hashed password", func(t *testing.T) {
		entryPassword := "password"
		hashedPassword, _ := GetHashedPassword(entryPassword)

		assert.NotEqual(t, entryPassword, hashedPassword)
	})

	t.Run("Hash password, should not be empty", func(t *testing.T) {
		entryPassword := "password"
		hashedPassword, _ := GetHashedPassword(entryPassword)

		assert.NotEqual(t, "", hashedPassword)
	})
}
