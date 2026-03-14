package handlers

import (
	"ad-bi-backend/config"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Campaign 广告活动
type Campaign struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Budget      float64 `json:"budget"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
}

// GetCampaigns 获取广告活动列表
func GetCampaigns(c *gin.Context) {
	db := config.GetDB()
	
	rows, err := db.Query(`
		SELECT id, name, description, status, budget, start_date, end_date 
		FROM ad_campaigns 
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

	campaigns := []Campaign{}
	for rows.Next() {
		var camp Campaign
		err := rows.Scan(&camp.ID, &camp.Name, &camp.Description, &camp.Status, 
			&camp.Budget, &camp.StartDate, &camp.EndDate)
		if err != nil {
			continue
		}
		campaigns = append(campaigns, camp)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    campaigns,
	})
}

// GetCampaignByID 根据 ID 获取广告活动
func GetCampaignByID(c *gin.Context) {
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
	var camp Campaign
	err = db.QueryRow(`
		SELECT id, name, description, status, budget, start_date, end_date 
		FROM ad_campaigns 
		WHERE id = ?
	`, id).Scan(&camp.ID, &camp.Name, &camp.Description, &camp.Status,
		&camp.Budget, &camp.StartDate, &camp.EndDate)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "活动不存在",
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
		"data":    camp,
	})
}

// CreateCampaign 创建广告活动
func CreateCampaign(c *gin.Context) {
	var camp Campaign
	if err := c.ShouldBindJSON(&camp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	db := config.GetDB()
	result, err := db.Exec(`
		INSERT INTO ad_campaigns (name, description, status, budget, start_date, end_date)
		VALUES (?, ?, ?, ?, ?, ?)
	`, camp.Name, camp.Description, camp.Status, camp.Budget, camp.StartDate, camp.EndDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建失败",
		})
		return
	}

	id, _ := result.LastInsertId()
	camp.ID = id

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data":    camp,
	})
}

// UpdateCampaign 更新广告活动
func UpdateCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的 ID",
		})
		return
	}

	var camp Campaign
	if err := c.ShouldBindJSON(&camp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	db := config.GetDB()
	_, err = db.Exec(`
		UPDATE ad_campaigns 
		SET name=?, description=?, status=?, budget=?, start_date=?, end_date=?
		WHERE id=?
	`, camp.Name, camp.Description, camp.Status, camp.Budget, camp.StartDate, camp.EndDate, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新失败",
		})
		return
	}

	camp.ID = id
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    camp,
	})
}

// DeleteCampaign 删除广告活动
func DeleteCampaign(c *gin.Context) {
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
	_, err = db.Exec("DELETE FROM ad_campaigns WHERE id = ?", id)
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
