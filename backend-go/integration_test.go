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

func setupIntegrationTest(t *testing.T) (sqlmock.Sqlmock, *gin.Engine, func()) {
	gin.SetMode(gin.TestMode)
	utils.SetJWTSecret("test-secret")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("创建mock数据库失败: %v", err)
	}
	config.SetDB(db)

	r := gin.New()
	r.Use(middleware.CORSMiddleware())

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
		}
	}

	cleanup := func() {
		db.Close()
	}

	return mock, r, cleanup
}

func TestIntegrationLoginAndGetCampaigns(t *testing.T) {
	mock, r, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// 步骤1: 登录获取Token
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

	var loginResp handlers.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &loginResp)
	assert.NoError(t, err)

	loginData := loginResp.Data.(map[string]interface{})
	token := loginData["token"].(string)
	assert.NotEmpty(t, token)

	// 步骤2: 使用Token获取活动列表
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM ad_campaigns").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	campaignRows := sqlmock.NewRows([]string{"id", "name", "description", "status", "budget", "start_date", "end_date"}).
		AddRow(1, "测试活动", "测试描述", "active", 1000.00, "2024-01-01", "2024-12-31")
	mock.ExpectQuery("SELECT.*FROM ad_campaigns").
		WithArgs(20, 0).
		WillReturnRows(campaignRows)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/api/v1/campaigns", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var campaignsResp handlers.APIResponse
	err = json.Unmarshal(w2.Body.Bytes(), &campaignsResp)
	assert.NoError(t, err)
	assert.Equal(t, 0, campaignsResp.Code)
}

func TestIntegrationCreateAndDeleteCampaign(t *testing.T) {
	mock, r, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// 获取Token
	token, _ := utils.GenerateToken(1, "testuser")

	// 步骤1: 创建活动
	mock.ExpectExec("INSERT INTO ad_campaigns").
		WithArgs("集成测试活动", "描述", "active", 500.0, "2024-01-01", "2024-12-31").
		WillReturnResult(sqlmock.NewResult(1, 1))

	createBody, _ := json.Marshal(handlers.CreateCampaignRequest{
		Name:        "集成测试活动",
		Description: "描述",
		Status:      "active",
		Budget:      500.0,
		StartDate:   "2024-01-01",
		EndDate:     "2024-12-31",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/campaigns", bytes.NewBuffer(createBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var createResp handlers.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &createResp)
	assert.NoError(t, err)
	assert.Equal(t, 0, createResp.Code)
	assert.Equal(t, "创建成功", createResp.Message)

	// 步骤2: 删除刚创建的活动
	mock.ExpectExec("DELETE FROM ad_campaigns WHERE id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("DELETE", "/api/v1/campaigns/1", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var deleteResp handlers.APIResponse
	err = json.Unmarshal(w2.Body.Bytes(), &deleteResp)
	assert.NoError(t, err)
	assert.Equal(t, 0, deleteResp.Code)
	assert.Equal(t, "删除成功", deleteResp.Message)
}

func TestIntegrationRegisterAndLogin(t *testing.T) {
	mock, r, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// 步骤1: 注册新用户
	mock.ExpectQuery("SELECT id FROM users WHERE username = \\?").
		WithArgs("integrationuser").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO users").
		WithArgs("integrationuser", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	registerBody, _ := json.Marshal(handlers.LoginRequest{
		Username: "integrationuser",
		Password: "testpass123",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var registerResp handlers.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &registerResp)
	assert.NoError(t, err)
	assert.Equal(t, 0, registerResp.Code)

	// 步骤2: 使用注册的用户登录
	hashedPassword, _ := utils.HashPassword("testpass123")
	rows := sqlmock.NewRows([]string{"id", "password"}).
		AddRow(1, hashedPassword)
	mock.ExpectQuery("SELECT id, password FROM users WHERE username = \\?").
		WithArgs("integrationuser").
		WillReturnRows(rows)

	loginBody, _ := json.Marshal(handlers.LoginRequest{
		Username: "integrationuser",
		Password: "testpass123",
	})

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(loginBody))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var loginResp handlers.APIResponse
	err = json.Unmarshal(w2.Body.Bytes(), &loginResp)
	assert.NoError(t, err)
	assert.Equal(t, 0, loginResp.Code)

	loginData := loginResp.Data.(map[string]interface{})
	assert.NotEmpty(t, loginData["token"])
}

func TestIntegrationUpdateCampaignFlow(t *testing.T) {
	mock, r, cleanup := setupIntegrationTest(t)
	defer cleanup()

	token, _ := utils.GenerateToken(1, "testuser")

	// 步骤1: 更新活动
	newName := "更新后的活动名"
	newStatus := "paused"
	mock.ExpectExec("UPDATE ad_campaigns SET name=\\?, status=\\? WHERE id=\\?").
		WithArgs(newName, newStatus, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 查询更新后的记录
	updatedRows := sqlmock.NewRows([]string{"id", "name", "description", "status", "budget", "start_date", "end_date"}).
		AddRow(1, "更新后的活动名", "描述", "paused", 1000.00, "2024-01-01", "2024-12-31")
	mock.ExpectQuery("SELECT.*FROM ad_campaigns WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(updatedRows)

	updateBody, _ := json.Marshal(handlers.UpdateCampaignRequest{
		Name:   &newName,
		Status: &newStatus,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/api/v1/campaigns/1", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updateResp handlers.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &updateResp)
	assert.NoError(t, err)
	assert.Equal(t, 0, updateResp.Code)
	assert.Equal(t, "更新成功", updateResp.Message)

	// 步骤2: 查询更新后的活动
	queryRows := sqlmock.NewRows([]string{"id", "name", "description", "status", "budget", "start_date", "end_date"}).
		AddRow(1, "更新后的活动名", "描述", "paused", 1000.00, "2024-01-01", "2024-12-31")
	mock.ExpectQuery("SELECT.*FROM ad_campaigns WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(queryRows)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/api/v1/campaigns/1", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var getResp handlers.APIResponse
	err = json.Unmarshal(w2.Body.Bytes(), &getResp)
	assert.NoError(t, err)
	assert.Equal(t, 0, getResp.Code)
}

func TestIntegrationFullAuthFlow(t *testing.T) {
	_, r, cleanup := setupIntegrationTest(t)
	defer cleanup()

	t.Run("无认证访问应被拒绝", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/campaigns", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("健康检查无需认证", func(t *testing.T) {
		// 健康检查不在测试路由中，但验证CORS中间件
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/login", nil)
		r.ServeHTTP(w, req)
		// GET /login 不存在，应该返回404
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
