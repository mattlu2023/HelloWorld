package handlers

import (
	"ad-bi-backend/config"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdUnit 广告单元
type AdUnit struct {
	ID           int64  `json:"id"`
	CampaignID   int64  `json:"campaign_id"`
	Name         string `json:"name"`
	AdType       string `json:"ad_type"`
	Placement    string `json:"placement"`
	CreativeURL  string `json:"creative_url"`
	LandingURL   string `json:"landing_url"`
	Status       string `json:"status"`
}

// GetAdUnits 获取广告单元列表
func GetAdUnits(c *gin.Context) {
	db := config.GetDB()
	
	rows, err := db.Query(`
		SELECT id, campaign_id, name, ad_type, placement, creative_url, landing_url, status
		FROM ad_units
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}
	defer rows.Close()

	units := []AdUnit{}
	for rows.Next() {
		var unit AdUnit
		err := rows.Scan(&unit.ID, &unit.CampaignID, &unit.Name, &unit.AdType,
			&unit.Placement, &unit.CreativeURL, &unit.LandingURL, &unit.Status)
		if err != nil {
			continue
		}
		units = append(units, unit)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    units,
	})
}

// GetAdUnitByID 根据 ID 获取广告单元
func GetAdUnitByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的 ID",
		})
		return
	}

	db := config.GetDB()
	var unit AdUnit
	err = db.QueryRow(`
		SELECT id, campaign_id, name, ad_type, placement, creative_url, landing_url, status
		FROM ad_units
		WHERE id = ?
	`, id).Scan(&unit.ID, &unit.CampaignID, &unit.Name, &unit.AdType,
		&unit.Placement, &unit.CreativeURL, &unit.LandingURL, &unit.Status)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "广告单元不存在",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    unit,
	})
}

// CreateAdUnit 创建广告单元
func CreateAdUnit(c *gin.Context) {
	var unit AdUnit
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	db := config.GetDB()
	result, err := db.Exec(`
		INSERT INTO ad_units (campaign_id, name, ad_type, placement, creative_url, landing_url, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, unit.CampaignID, unit.Name, unit.AdType, unit.Placement, 
		unit.CreativeURL, unit.LandingURL, unit.Status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建失败",
		})
		return
	}

	id, _ := result.LastInsertId()
	unit.ID = id

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    unit,
	})
}

// UpdateAdUnit 更新广告单元
func UpdateAdUnit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的 ID",
		})
		return
	}

	var unit AdUnit
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	db := config.GetDB()
	_, err = db.Exec(`
		UPDATE ad_units 
		SET campaign_id=?, name=?, ad_type=?, placement=?, creative_url=?, landing_url=?, status=?
		WHERE id=?
	`, unit.CampaignID, unit.Name, unit.AdType, unit.Placement,
		unit.CreativeURL, unit.LandingURL, unit.Status, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新失败",
		})
		return
	}

	unit.ID = id
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    unit,
	})
}

// DeleteAdUnit 删除广告单元
func DeleteAdUnit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的 ID",
		})
		return
	}

	db := config.GetDB()
	_, err = db.Exec("DELETE FROM ad_units WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}
