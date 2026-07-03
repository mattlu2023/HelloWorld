package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetJWTSecret(t *testing.T) {
	secret := "test-secret-key"
	SetJWTSecret(secret)

	// 验证密钥已设置（通过生成token间接验证）
	token, err := GenerateToken(1, "testuser")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateToken(t *testing.T) {
	SetJWTSecret("test-secret-key")

	tests := []struct {
		name     string
		userID   int64
		username string
		wantErr  bool
	}{
		{
			name:     "正常生成token",
			userID:   1,
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "用户ID为0",
			userID:   0,
			username: "testuser",
			wantErr:  false,
		},
		{
			name:     "用户名为空",
			userID:   1,
			username: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.username)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	SetJWTSecret("test-secret-key")

	t.Run("正常解析token", func(t *testing.T) {
		userID := int64(123)
		username := "testuser"
		token, err := GenerateToken(userID, username)
		assert.NoError(t, err)

		claims, err := ParseToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)

		// 验证claims中的用户信息
		// 注意：JSON数字解析后为float64
		assert.Equal(t, float64(userID), claims["user_id"])
		assert.Equal(t, username, claims["username"])
	})

	t.Run("无效的token", func(t *testing.T) {
		claims, err := ParseToken("invalid-token-string")
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("空token", func(t *testing.T) {
		claims, err := ParseToken("")
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("使用不同密钥签名的token", func(t *testing.T) {
		SetJWTSecret("secret-1")
		token, err := GenerateToken(1, "user1")
		assert.NoError(t, err)

		SetJWTSecret("secret-2")
		claims, err := ParseToken(token)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestTokenExpiration(t *testing.T) {
	SetJWTSecret("test-secret-key")

	t.Run("token应包含过期时间", func(t *testing.T) {
		token, err := GenerateToken(1, "testuser")
		assert.NoError(t, err)

		claims, err := ParseToken(token)
		assert.NoError(t, err)

		exp, ok := claims["exp"]
		assert.True(t, ok)
		expFloat, ok := exp.(float64)
		assert.True(t, ok)

		// 验证过期时间是24小时后
		expTime := time.Unix(int64(expFloat), 0)
		expectedExp := time.Now().Add(24 * time.Hour)
		assert.WithinDuration(t, expectedExp, expTime, 5*time.Second)
	})
}
