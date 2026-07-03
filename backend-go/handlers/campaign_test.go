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

func setupCampaignTest(t *testing.T) (sqlmock.Sqlmock, func()) {
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

func TestGetCampaigns(t *testing.T) {
	t.Run("获取活动列表成功", func(t *testing.T) {
		mock, cleanup := setupCampaignTest(t)
		defer cleanup()

		// 查询总数
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM ad_campaigns").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		// 查询列表
		rows := sqlmock.NewRows([]string{"id", "name", "description", "status", "budget", "start_date", "end_date"}).
			AddRow(1, "活动A", "描述A", "active", 1000.00, "2024-01-01", "2024-12-31").
			AddRow(2, "活动B", "描述B", "paused", 2000.00, "2024-02-01", "2024-11-30")
		mock.ExpectQuery("SELECT id, name, description, status, budget, start_date, end_date FROM ad_campaigns").
			WithArgs(20, 0).
			WillReturnRows(rows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/campaigns?page=1&page_size=20", nil)

		GetCampaigns(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 0, resp.Code)
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		mock, cleanup := setupCampaignTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM ad_campaigns").
			WillReturnError(sql.ErrConnDone)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/campaigns", nil)

		GetCampaigns(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetCampaignByID(t *testing.T) {
	t.Run("获取活动详情成功", func(t *testing.T) {
		mock, cleanup := setupCampaignTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"id", "name", "description", "status", "budget", "start_date", "end_date"}).
			AddRow(1, "活动A", "描述A", "active", 1000.00, "2024-01-01", "2024-12-31")
		mock.ExpectQuery("SELECT id, name, description, status, budget, start_date, end_date FROM ad_campaigns WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(rows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/campaigns/1", nil)

		GetCampaignByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("活动不存在", func(t *testing.T) {
		mock, cleanup := setupCampaignTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT id, name, description, status, budget, start_date, end_date FROM ad_campaigns WHERE id = \\?").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request = httptest.NewRequest("GET", "/campaigns/999", nil)

		GetCampaignByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("无效的ID", func(t *testing.T) {
		_, cleanup := setupCampaignTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request = httptest.NewRequest("GET", "/campaigns/abc", nil)

		GetCampaignByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCreateCampaign(t *testing.T) {
	t.Run("创建活动成功", func(t *testing.T) {
		mock, cleanup := setupCampaignTest(t)
		defer cleanup()

		mock.ExpectExec("INSERT INTO ad_campaigns").
			WithArgs("新活动", "描述", "active", 1000.0, "2024-01-01", "2024-12-31").
			WillReturnResult(sqlmock.NewResult(1, 1))

		req := CreateCampaignRequest{
			Name:        "新活动",
			Description: "描述",
			Status:      "active",
			Budget:      1000.0,
			StartDate:   "2024-01-01",
			EndDate:     "2024-12-31",
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/campaigns", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		CreateCampaign(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 0, resp.Code)
		assert.Equal(t, "创建成功", resp.Message)
	})

	t.Run("请求参数错误-缺少必填字段", func(t *testing.T) {
		_, cleanup := setupCampaignTest(t)
		defer cleanup()

		body, _ := json.Marshal(map[string]interface{}{
			"description": "描述",
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/campaigns", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		CreateCampaign(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("无效的状态值", func(t *testing.T) {
		_, cleanup := setupCampaignTest(t)
		defer cleanup()

		body, _ := json.Marshal(CreateCampaignRequest{
			Name:      "活动",
			Status:    "invalid",
			Budget:    1000.0,
			StartDate: "2024-01-01",
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/campaigns", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		CreateCampaign(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUpdateCampaign(t *testing.T) {
	t.Run("更新活动成功", func(t *testing.T) {
		mock, cleanup := setupCampaignTest(t)
		defer cleanup()

		name := "更新后的名称"
		status := "paused"
		mock.ExpectExec("UPDATE ad_campaigns SET name=\\?, status=\\? WHERE id=\\?").
			WithArgs(name, status, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// 查询更新后的记录
		rows := sqlmock.NewRows([]string{"id", "name", "description", "status", "budget", "start_date", "end_date"}).
			AddRow(1, "更新后的名称", "描述", "paused", 1000.00, "2024-01-01", "2024-12-31")
		mock.ExpectQuery("SELECT id, name, description, status, budget, start_date, end_date FROM ad_campaigns WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(rows)

		req := UpdateCampaignRequest{
			Name:   &name,
			Status: &status,
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("PUT", "/campaigns/1", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		UpdateCampaign(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("活动不存在", func(t *testing.T) {
		mock, cleanup := setupCampaignTest(t)
		defer cleanup()

		name := "新名称"
		mock.ExpectExec("UPDATE ad_campaigns SET name=\\? WHERE id=\\?").
			WithArgs(name, 999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		req := UpdateCampaignRequest{
			Name: &name,
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request = httptest.NewRequest("PUT", "/campaigns/999", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		UpdateCampaign(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("未提供任何更新字段", func(t *testing.T) {
		_, cleanup := setupCampaignTest(t)
		defer cleanup()

		body, _ := json.Marshal(UpdateCampaignRequest{})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("PUT", "/campaigns/1", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		UpdateCampaign(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteCampaign(t *testing.T) {
	t.Run("删除活动成功", func(t *testing.T) {
		mock, cleanup := setupCampaignTest(t)
		defer cleanup()

		mock.ExpectExec("DELETE FROM ad_campaigns WHERE id = \\?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("DELETE", "/campaigns/1", nil)

		DeleteCampaign(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 0, resp.Code)
		assert.Equal(t, "删除成功", resp.Message)
	})

	t.Run("活动不存在", func(t *testing.T) {
		mock, cleanup := setupCampaignTest(t)
		defer cleanup()

		mock.ExpectExec("DELETE FROM ad_campaigns WHERE id = \\?").
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request = httptest.NewRequest("DELETE", "/campaigns/999", nil)

		DeleteCampaign(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("无效的ID", func(t *testing.T) {
		_, cleanup := setupCampaignTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request = httptest.NewRequest("DELETE", "/campaigns/abc", nil)

		DeleteCampaign(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
