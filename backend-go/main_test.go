package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ad-bi-backend/config"
	"ad-bi-backend/handlers"
	"ad-bi-backend/middleware"
	"ad-bi-backend/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAPITest(t *testing.T) (sqlmock.Sqlmock, *gin.Engine, func()) {
	gin.SetMode(gin.TestMode)
	utils.SetJWTSecret("test-secret")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("创建mock数据库失败: %v", err)
	}
	config.SetDB(db)

	r := gin.New()
	r.Use(middleware.CORSMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "广告 BI 系统运行正常",
		})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/login", handlers.Login)
		api.POST("/register", handlers.Register)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/campaigns", handlers.GetCampaigns)
			auth.GET("/campaigns/:id", handlers.GetCampaignByID)
			auth.POST("/campaigns", handlers.CreateCampaign)
			auth.PUT("/campaigns/:id", handlers.UpdateCampaign)
			auth.DELETE("/campaigns/:id", handlers.DeleteCampaign)

			auth.GET("/ad-units", handlers.GetAdUnits)
			auth.POST("/ad-units", handlers.CreateAdUnit)
			auth.DELETE("/ad-units/:id", handlers.DeleteAdUnit)

			auth.GET("/stats/overview", handlers.GetOverview)
			auth.GET("/stats/daily-trend", handlers.GetDailyTrend)
			auth.GET("/stats/funnel", handlers.GetConversionFunnel)
		}
	}

	cleanup := func() {
		db.Close()
	}

	return mock, r, cleanup
}

func TestHealthEndpoint(t *testing.T) {
	_, r, cleanup := setupAPITest(t)
	defer cleanup()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp["status"])
}

func TestCORSMiddleware(t *testing.T) {
	_, r, cleanup := setupAPITest(t)
	defer cleanup()

	t.Run("OPTIONS预检请求", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("OPTIONS", "/health", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("带Origin的GET请求", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/health", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("无Origin的请求", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/health", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestAuthMiddleware(t *testing.T) {
	_, r, cleanup := setupAPITest(t)
	defer cleanup()

	t.Run("无Token访问受保护接口", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/campaigns", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("无效Token访问受保护接口", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/campaigns", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("有效Token访问受保护接口", func(t *testing.T) {
		mock, _, _ := setupAPITest(t)
		defer func() {
			config.SetDB(nil)
		}()

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM ad_campaigns").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectQuery("SELECT.*FROM ad_campaigns").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "status", "budget", "start_date", "end_date"}))

		token, err := utils.GenerateToken(1, "testuser")
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/campaigns", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAPILoginFlow(t *testing.T) {
	mock, r, cleanup := setupAPITest(t)
	defer cleanup()

	t.Run("API登录成功", func(t *testing.T) {
		hashedPassword, _ := utils.HashPassword("password123")

		rows := sqlmock.NewRows([]string{"id", "password"}).
			AddRow(1, hashedPassword)
		mock.ExpectQuery("SELECT id, password FROM users WHERE username = \\?").
			WithArgs("testuser").
			WillReturnRows(rows)

		body, _ := json.Marshal(handlers.LoginRequest{
			Username: "testuser",
			Password: "password123",
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp handlers.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 0, resp.Code)

		data := resp.Data.(map[string]interface{})
		assert.NotEmpty(t, data["token"])
	})

	t.Run("API登录失败-密码错误", func(t *testing.T) {
		hashedPassword, _ := utils.HashPassword("correctpassword")

		rows := sqlmock.NewRows([]string{"id", "password"}).
			AddRow(1, hashedPassword)
		mock.ExpectQuery("SELECT id, password FROM users WHERE username = \\?").
			WithArgs("testuser").
			WillReturnRows(rows)

		body, _ := json.Marshal(handlers.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAPIRegisterFlow(t *testing.T) {
	mock, r, cleanup := setupAPITest(t)
	defer cleanup()

	t.Run("API注册成功", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM users WHERE username = \\?").
			WithArgs("newuser").
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec("INSERT INTO users").
			WithArgs("newuser", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		body, _ := json.Marshal(handlers.LoginRequest{
			Username: "newuser",
			Password: "password123",
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
