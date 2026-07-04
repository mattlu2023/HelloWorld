package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("正常加密密码", func(t *testing.T) {
		password := "mySecretPassword123"
		hashed, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)
		assert.NotEqual(t, password, hashed)
	})

	t.Run("相同密码每次加密结果不同", func(t *testing.T) {
		password := "testPassword"
		hash1, err := HashPassword(password)
		assert.NoError(t, err)

		hash2, err := HashPassword(password)
		assert.NoError(t, err)

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("空密码", func(t *testing.T) {
		hashed, err := HashPassword("")
		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)
	})

	t.Run("长密码", func(t *testing.T) {
		password := strings.Repeat("a", 72) // bcrypt最大长度72字节
		hashed, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hashed)
	})
}

func TestVerifyPassword(t *testing.T) {
	t.Run("正确密码验证成功", func(t *testing.T) {
		password := "correctPassword"
		hashed, err := HashPassword(password)
		assert.NoError(t, err)

		assert.True(t, VerifyPassword(hashed, password))
	})

	t.Run("错误密码验证失败", func(t *testing.T) {
		password := "correctPassword"
		hashed, err := HashPassword(password)
		assert.NoError(t, err)

		assert.False(t, VerifyPassword(hashed, "wrongPassword"))
	})

	t.Run("空密码与有效哈希", func(t *testing.T) {
		password := "somePassword"
		hashed, err := HashPassword(password)
		assert.NoError(t, err)

		assert.False(t, VerifyPassword(hashed, ""))
	})

	t.Run("无效的哈希字符串", func(t *testing.T) {
		assert.False(t, VerifyPassword("invalid-hash", "password"))
	})

	t.Run("空哈希字符串", func(t *testing.T) {
		assert.False(t, VerifyPassword("", "password"))
	})

	t.Run("大小写敏感", func(t *testing.T) {
		password := "Password123"
		hashed, err := HashPassword(password)
		assert.NoError(t, err)

		assert.False(t, VerifyPassword(hashed, "password123"))
	})
}
