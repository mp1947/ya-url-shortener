package auth

import (
	"testing"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func createToken() (string, error) {
	userID := uuid.New()
	return CreateToken(userID)
}

func TestCreateToken(t *testing.T) {
	t.Run("create_token", func(t *testing.T) {
		token, err := createToken()
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestValidate(t *testing.T) {

	t.Run("validate_correct_token", func(t *testing.T) {
		token, err := createToken()
		assert.NoError(t, err)
		isOk, userID := Validate(token)
		assert.True(t, isOk)
		assert.NotNil(t, userID)
	})

	t.Run("validate_incorrect_token", func(t *testing.T) {
		token, err := createToken()
		assert.NoError(t, err)
		isOk, userID := Validate(token + "aaa")
		assert.False(t, isOk)
		assert.Equal(t, uuid.Nil, userID)
	})
}

func BenchmarkCreateToken(b *testing.B) {

	for i := 0; i < b.N; i++ {
		_, _ = createToken()
	}
	b.ReportAllocs()
}
