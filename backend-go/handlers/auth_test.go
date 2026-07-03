package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ad-bi-backend/config"
	"ad-bi-backend/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthTest(t *testing.T) (sqlmock.Sqlmock, func()) {
	gin.SetMode(gin.TestMode)
	utils.SetJWTSecret("test-secret")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("创建mock数据库失败: %v", err)
	}

	config.SetDB(db)

	cleanup := func() {
		db.Close()
	}

	return mock, cleanup
}

func TestLogin(t *testing.T) {
	t.Run("登录成功", func(t *testing.T) {
		mock, cleanup := setupAuthTest(t)
		defer cleanup()

		hashedPassword, _ := utils.HashPassword("password123")

		rows := sqlmock.NewRows([]string{"id", "password"}).
			AddRow(1, hashedPassword)
		mock.ExpectQuery("SELECT id, password FROM users WHERE username = \\?").
			WithArgs("testuser").
			WillReturnRows(rows)

		body, _ := json.Marshal(LoginRequest{
			Username: "testuser",
			Password: "password123",
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		Login(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 0, resp.Code)
		assert.Equal(t, "登录成功", resp.Message)

		data := resp.Data.(map[string]interface{})
		assert.Equal(t, float64(1), data["user_id"])
		assert.Equal(t, "testuser", data["username"])
		assert.NotEmpty(t, data["token"])
	})

	t.Run("用户不存在", func(t *testing.T) {
		mock, cleanup := setupAuthTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT id, password FROM users WHERE username = \\?").
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		body, _ := json.Marshal(LoginRequest{
			Username: "nonexistent",
			Password: "password123",
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var resp APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.Code)
		assert.Equal(t, "用户名或密码错误", resp.Message)
	})

	t.Run("密码错误", func(t *testing.T) {
		mock, cleanup := setupAuthTest(t)
		defer cleanup()

		hashedPassword, _ := utils.HashPassword("correctpassword")

		rows := sqlmock.NewRows([]string{"id", "password"}).
			AddRow(1, hashedPassword)
		mock.ExpectQuery("SELECT id, password FROM users WHERE username = \\?").
			WithArgs("testuser").
			WillReturnRows(rows)

		body, _ := json.Marshal(LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var resp APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.Code)
		assert.Equal(t, "用户名或密码错误", resp.Message)
	})

	t.Run("请求参数错误-缺少用户名", func(t *testing.T) {
		_, cleanup := setupAuthTest(t)
		defer cleanup()

		body, _ := json.Marshal(map[string]string{
			"password": "password123",
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRegister(t *testing.T) {
	t.Run("注册成功", func(t *testing.T) {
		mock, cleanup := setupAuthTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT id FROM users WHERE username = \\?").
			WithArgs("newuser").
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec("INSERT INTO users \\(username, password\\) VALUES \\(\\?, \\?\\)").
			WithArgs("newuser", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		body, _ := json.Marshal(LoginRequest{
			Username: "newuser",
			Password: "password123",
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		Register(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 0, resp.Code)
		assert.Equal(t, "注册成功", resp.Message)
	})

	t.Run("用户名已存在", func(t *testing.T) {
		mock, cleanup := setupAuthTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("SELECT id FROM users WHERE username = \\?").
			WithArgs("existinguser").
			WillReturnRows(rows)

		body, _ := json.Marshal(LoginRequest{
			Username: "existinguser",
			Password: "password123",
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		Register(c)

		assert.Equal(t, http.StatusConflict, w.Code)

		var resp APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 409, resp.Code)
		assert.Equal(t, "用户名已存在，请选择其他用户名", resp.Message)
	})

	t.Run("请求参数错误", func(t *testing.T) {
		_, cleanup := setupAuthTest(t)
		defer cleanup()

		body := []byte(`{"invalid json`)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
