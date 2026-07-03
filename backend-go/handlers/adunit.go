package handlers

import (
	"ad-bi-backend/config"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

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

type CreateAdUnitRequest struct {
	CampaignID   int64  `json:"campaign_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	AdType       string `json:"ad_type" binding:"required"`
	Placement    string `json:"placement"`
	CreativeURL  string `json:"creative_url"`
	LandingURL   string `json:"landing_url"`
	Status       string `json:"status" binding:"required,oneof=active paused"`
}

type UpdateAdUnitRequest struct {
	CampaignID   *int64  `json:"campaign_id"`
	Name         *string `json:"name"`
	AdType       *string `json:"ad_type"`
	Placement    *string `json:"placement"`
	CreativeURL  *string `json:"creative_url"`
	LandingURL   *string `json:"landing_url"`
	Status       *string `json:"status"`
}

func GetAdUnits(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	db := config.GetDB()

	// 查询总数
	var total int64
	err := db.QueryRow("SELECT COUNT(*) FROM ad_units").Scan(&total)
	if err != nil {
		log.Printf("查询广告单元总数失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}

	// 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	rows, err := db.Query(`
		SELECT id, campaign_id, name, ad_type, placement, creative_url, landing_url, status
		FROM ad_units
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, pageSize, offset)
	if err != nil {
		log.Printf("查询广告单元列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}
	defer rows.Close()

	units := []AdUnit{}
	for rows.Next() {
		var unit AdUnit
		if err := rows.Scan(&unit.ID, &unit.CampaignID, &unit.Name, &unit.AdType,
			&unit.Placement, &unit.CreativeURL, &unit.LandingURL, &unit.Status); err != nil {
			log.Printf("扫描广告单元数据失败: %v", err)
			continue
		}
		units = append(units, unit)
	}

	if err := rows.Err(); err != nil {
		log.Printf("遍历广告单元数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
		Data: PaginationResponse{
			List:       units,
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func GetAdUnitByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的 ID",
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
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "广告单元不存在",
		})
		return
	}
	if err != nil {
		log.Printf("查询广告单元失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
		Data:    unit,
	})
}

func CreateAdUnit(c *gin.Context) {
	var req CreateAdUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "请求参数错误：活动ID、名称、广告类型和状态为必填项，状态只能为 active/paused",
		})
		return
	}

	db := config.GetDB()
	result, err := db.Exec(`
		INSERT INTO ad_units (campaign_id, name, ad_type, placement, creative_url, landing_url, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, req.CampaignID, req.Name, req.AdType, req.Placement,
		req.CreativeURL, req.LandingURL, req.Status)

	if err != nil {
		log.Printf("创建广告单元失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "创建广告单元失败，请稍后重试",
		})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "创建成功",
		Data: AdUnit{
			ID:           id,
			CampaignID:   req.CampaignID,
			Name:         req.Name,
			AdType:       req.AdType,
			Placement:    req.Placement,
			CreativeURL:  req.CreativeURL,
			LandingURL:   req.LandingURL,
			Status:       req.Status,
		},
	})
}

func UpdateAdUnit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的 ID",
		})
		return
	}

	var req UpdateAdUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	db := config.GetDB()

	// 构建动态UPDATE语句，只更新非nil字段
	var updateFields []string
	var args []interface{}

	if req.CampaignID != nil {
		updateFields = append(updateFields, "campaign_id=?")
		args = append(args, *req.CampaignID)
	}
	if req.Name != nil {
		updateFields = append(updateFields, "name=?")
		args = append(args, *req.Name)
	}
	if req.AdType != nil {
		updateFields = append(updateFields, "ad_type=?")
		args = append(args, *req.AdType)
	}
	if req.Placement != nil {
		updateFields = append(updateFields, "placement=?")
		args = append(args, *req.Placement)
	}
	if req.CreativeURL != nil {
		updateFields = append(updateFields, "creative_url=?")
		args = append(args, *req.CreativeURL)
	}
	if req.LandingURL != nil {
		updateFields = append(updateFields, "landing_url=?")
		args = append(args, *req.LandingURL)
	}
	if req.Status != nil {
		updateFields = append(updateFields, "status=?")
		args = append(args, *req.Status)
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "未提供任何更新字段",
		})
		return
	}

	// 检查状态值是否合法
	if req.Status != nil && *req.Status != "active" && *req.Status != "paused" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的状态值，可选值: active, paused",
		})
		return
	}

	sql := "UPDATE ad_units SET " + strings.Join(updateFields, ", ") + " WHERE id=?"
	args = append(args, id)

	result, err := db.Exec(sql, args...)
	if err != nil {
		log.Printf("更新广告单元失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "更新失败",
		})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取更新影响行数失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "更新失败",
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "广告单元不存在",
		})
		return
	}

	// 查询更新后的完整记录
	var updatedUnit AdUnit
	err = db.QueryRow(`
		SELECT id, campaign_id, name, ad_type, placement, creative_url, landing_url, status
		FROM ad_units
		WHERE id = ?
	`, id).Scan(&updatedUnit.ID, &updatedUnit.CampaignID, &updatedUnit.Name, &updatedUnit.AdType,
		&updatedUnit.Placement, &updatedUnit.CreativeURL, &updatedUnit.LandingURL, &updatedUnit.Status)

	if err != nil {
		log.Printf("查询更新后广告单元失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询更新结果失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "更新成功",
		Data:    updatedUnit,
	})
}

func DeleteAdUnit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的 ID",
		})
		return
	}

	db := config.GetDB()
	result, err := db.Exec("DELETE FROM ad_units WHERE id = ?", id)
	if err != nil {
		log.Printf("删除广告单元失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "删除失败",
		})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取删除影响行数失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "删除失败",
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "广告单元不存在",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "删除成功",
	})
}