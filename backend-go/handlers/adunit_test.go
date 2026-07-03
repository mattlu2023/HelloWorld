package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"ad-bi-backend/config"
	"ad-bi-backend/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAdUnitTest(t *testing.T) (sqlmock.Sqlmock, func()) {
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

func TestGetAdUnits(t *testing.T) {
	t.Run("获取广告单元列表成功", func(t *testing.T) {
		mock, cleanup := setupAdUnitTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM ad_units").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		rows := sqlmock.NewRows([]string{"id", "campaign_id", "name", "ad_type", "placement", "creative_url", "landing_url", "status"}).
			AddRow(1, 1, "广告单元A", "banner", "top", "http://example.com/a.jpg", "http://example.com", "active").
			AddRow(2, 1, "广告单元B", "video", "sidebar", "http://example.com/b.mp4", "http://example.com", "paused")
		mock.ExpectQuery("SELECT id, campaign_id, name, ad_type, placement, creative_url, landing_url, status FROM ad_units").
			WithArgs(20, 0).
			WillReturnRows(rows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/ad-units?page=1&page_size=20", nil)

		GetAdUnits(c)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		mock, cleanup := setupAdUnitTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM ad_units").
			WillReturnError(sql.ErrConnDone)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/ad-units", nil)

		GetAdUnits(c)

		assert.Equal(t, 500, w.Code)
	})
}

func TestGetAdUnitByID(t *testing.T) {
	t.Run("获取广告单元详情成功", func(t *testing.T) {
		mock, cleanup := setupAdUnitTest(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"id", "campaign_id", "name", "ad_type", "placement", "creative_url", "landing_url", "status"}).
			AddRow(1, 1, "广告单元A", "banner", "top", "http://example.com/a.jpg", "http://example.com", "active")
		mock.ExpectQuery("SELECT id, campaign_id, name, ad_type, placement, creative_url, landing_url, status FROM ad_units WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(rows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("GET", "/ad-units/1", nil)

		GetAdUnitByID(c)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("广告单元不存在", func(t *testing.T) {
		mock, cleanup := setupAdUnitTest(t)
		defer cleanup()

		mock.ExpectQuery("SELECT id, campaign_id, name, ad_type, placement, creative_url, landing_url, status FROM ad_units WHERE id = \\?").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request = httptest.NewRequest("GET", "/ad-units/999", nil)

		GetAdUnitByID(c)

		assert.Equal(t, 404, w.Code)
	})
}

func TestCreateAdUnit(t *testing.T) {
	t.Run("创建广告单元成功", func(t *testing.T) {
		mock, cleanup := setupAdUnitTest(t)
		defer cleanup()

		mock.ExpectExec("INSERT INTO ad_units").
			WithArgs(1, "新广告", "banner", "top", "http://example.com/a.jpg", "http://example.com", "active").
			WillReturnResult(sqlmock.NewResult(1, 1))

		req := CreateAdUnitRequest{
			CampaignID:  1,
			Name:        "新广告",
			AdType:      "banner",
			Placement:   "top",
			CreativeURL: "http://example.com/a.jpg",
			LandingURL:  "http://example.com",
			Status:      "active",
		}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/ad-units", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		CreateAdUnit(c)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("缺少必填字段", func(t *testing.T) {
		_, cleanup := setupAdUnitTest(t)
		defer cleanup()

		body, _ := json.Marshal(map[string]interface{}{
			"name": "广告",
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/ad-units", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		CreateAdUnit(c)

		assert.Equal(t, 400, w.Code)
	})
}

func TestDeleteAdUnit(t *testing.T) {
	t.Run("删除广告单元成功", func(t *testing.T) {
		mock, cleanup := setupAdUnitTest(t)
		defer cleanup()

		mock.ExpectExec("DELETE FROM ad_units WHERE id = \\?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("DELETE", "/ad-units/1", nil)

		DeleteAdUnit(c)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("广告单元不存在", func(t *testing.T) {
		mock, cleanup := setupAdUnitTest(t)
		defer cleanup()

		mock.ExpectExec("DELETE FROM ad_units WHERE id = \\?").
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request = httptest.NewRequest("DELETE", "/ad-units/999", nil)

		DeleteAdUnit(c)

		assert.Equal(t, 404, w.Code)
	})
}
