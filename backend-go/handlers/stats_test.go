package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"ad-bi-backend/config"
	"ad-bi-backend/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupStatsTest(t *testing.T) (sqlmock.Sqlmock, func()) {
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

func TestGetOverview(t *testing.T) {
	t.Run("获取数据概览成功", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{
			"total_impressions", "total_clicks", "total_conversions",
			"total_cost", "total_revenue", "average_ctr", "average_conversion_rate",
		}).AddRow(10000, 500, 50, 1000.00, 5000.00, 5.0, 10.0)

		mock.ExpectQuery("SELECT.*FROM ad_stats_daily").
			WillReturnRows(rows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/stats/overview", nil)

		GetOverview(c)

		assert.Equal(t, 200, w.Code)

		var resp APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 0, resp.Code)
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT.*FROM ad_stats_daily").
			WillReturnError(sqlmock.ErrCancelled)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/stats/overview", nil)

		GetOverview(c)

		assert.Equal(t, 500, w.Code)
	})
}

func TestGetDailyTrend(t *testing.T) {
	t.Run("获取每日趋势成功", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{
			"stat_date", "impressions", "clicks", "conversions", "cost", "revenue", "ctr",
		}).
			AddRow("2024-01-01", 1000, 50, 5, 100.00, 500.00, 5.0).
			AddRow("2024-01-02", 1200, 60, 6, 120.00, 600.00, 5.0)

		mock.ExpectQuery("SELECT.*FROM ad_stats_daily.*GROUP BY stat_date").
			WillReturnRows(rows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/stats/daily-trend", nil)

		GetDailyTrend(c)

		assert.Equal(t, 200, w.Code)

		var resp APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 0, resp.Code)
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT.*FROM ad_stats_daily").
			WillReturnError(sqlmock.ErrCancelled)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/stats/daily-trend", nil)

		GetDailyTrend(c)

		assert.Equal(t, 500, w.Code)
	})
}

func TestGetConversionFunnel(t *testing.T) {
	t.Run("获取转化漏斗成功", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"step_name", "step_order", "user_count"}).
			AddRow("曝光", 1, 10000).
			AddRow("点击", 2, 500).
			AddRow("转化", 3, 50)

		mock.ExpectQuery("SELECT step_name, step_order, user_count FROM conversion_funnel").
			WillReturnRows(rows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/stats/funnel", nil)

		GetConversionFunnel(c)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT step_name, step_order, user_count FROM conversion_funnel").
			WillReturnError(sqlmock.ErrCancelled)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/stats/funnel", nil)

		GetConversionFunnel(c)

		assert.Equal(t, 500, w.Code)
	})
}

func TestGetCampaignStats(t *testing.T) {
	t.Run("获取活动统计成功", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{
			"stat_date", "impressions", "clicks", "conversions", "cost", "revenue",
		}).
			AddRow("2024-01-01", 1000, 50, 5, 100.00, 500.00)

		mock.ExpectQuery("SELECT.*FROM ad_stats_daily WHERE campaign_id = \\?").
			WithArgs(1).
			WillReturnRows(rows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/stats/campaign/1", nil)

		GetCampaignStats(c)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("无效的ID", func(t *testing.T) {
		_, cleanup := setupStatsTest(t)
		defer cleanup()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request = httptest.NewRequest("GET", "/stats/campaign/abc", nil)

		GetCampaignStats(c)

		assert.Equal(t, 400, w.Code)
	})
}

func TestGetUserActions(t *testing.T) {
	t.Run("获取用户行为列表成功", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM user_actions").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		rows := sqlmock.NewRows([]string{"id", "user_id", "action", "page", "created_at"}).
			AddRow(1, 1, "click", "home", "2024-01-01 10:00:00").
			AddRow(2, 1, "view", "dashboard", "2024-01-01 11:00:00")

		mock.ExpectQuery("SELECT id, user_id, action, page, created_at FROM user_actions").
			WithArgs(20, 0).
			WillReturnRows(rows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/user-actions?page=1&page_size=20", nil)

		GetUserActions(c)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM user_actions").
			WillReturnError(sqlmock.ErrCancelled)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/user-actions", nil)

		GetUserActions(c)

		assert.Equal(t, 500, w.Code)
	})
}

func TestGetUserActionAnalysis(t *testing.T) {
	t.Run("获取用户行为分析成功", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"action", "count"}).
			AddRow("click", 100).
			AddRow("view", 200)

		mock.ExpectQuery("SELECT action, COUNT\\(\\*\\) as count FROM user_actions GROUP BY action").
			WillReturnRows(rows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/user-actions/analysis", nil)

		GetUserActionAnalysis(c)

		assert.Equal(t, 200, w.Code)

		var resp APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 0, resp.Code)
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT action, COUNT\\(\\*\\) as count FROM user_actions").
			WillReturnError(sqlmock.ErrCancelled)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/user-actions/analysis", nil)

		GetUserActionAnalysis(c)

		assert.Equal(t, 500, w.Code)
	})
}

func TestExportReport(t *testing.T) {
	t.Run("导出报表成功", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{
			"stat_date", "campaign_id", "impressions", "clicks",
			"conversions", "cost", "revenue", "ctr", "conversion_rate",
		}).
			AddRow("2024-01-01", 1, 1000, 50, 5, 100.00, 500.00, 5.0, 10.0)

		mock.ExpectQuery("SELECT.*FROM ad_stats_daily").
			WillReturnRows(rows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/reports/export", nil)

		ExportReport(c)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")
		assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
	})

	t.Run("带日期范围导出", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{
			"stat_date", "campaign_id", "impressions", "clicks",
			"conversions", "cost", "revenue", "ctr", "conversion_rate",
		}).
			AddRow("2024-01-01", 1, 1000, 50, 5, 100.00, 500.00, 5.0, 10.0)

		mock.ExpectQuery("SELECT.*FROM ad_stats_daily WHERE stat_date BETWEEN").
			WithArgs("2024-01-01", "2024-01-31").
			WillReturnRows(rows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/reports/export?start_date=2024-01-01&end_date=2024-01-31", nil)

		ExportReport(c)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		mock, cleanup := setupStatsTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT.*FROM ad_stats_daily").
			WillReturnError(sqlmock.ErrCancelled)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/reports/export", nil)

		ExportReport(c)

		assert.Equal(t, 500, w.Code)
	})
}
