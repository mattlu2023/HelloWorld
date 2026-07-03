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

type Campaign struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Budget      float64 `json:"budget"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
}

type CreateCampaignRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Status      string  `json:"status" binding:"required,oneof=active paused ended"`
	Budget      float64 `json:"budget" binding:"required,min=0"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     string  `json:"end_date"`
}

type UpdateCampaignRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Status      *string  `json:"status"`
	Budget      *float64 `json:"budget"`
	StartDate   *string  `json:"start_date"`
	EndDate     *string  `json:"end_date"`
}

func GetCampaigns(c *gin.Context) {
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
	err := db.QueryRow("SELECT COUNT(*) FROM ad_campaigns").Scan(&total)
	if err != nil {
		log.Printf("查询广告活动总数失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}

	// 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	rows, err := db.Query(`
		SELECT id, name, description, status, budget, start_date, end_date 
		FROM ad_campaigns 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, pageSize, offset)
	if err != nil {
		log.Printf("查询广告活动列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}
	defer rows.Close()

	campaigns := []Campaign{}
	for rows.Next() {
		var camp Campaign
		if err := rows.Scan(&camp.ID, &camp.Name, &camp.Description, &camp.Status,
			&camp.Budget, &camp.StartDate, &camp.EndDate); err != nil {
			log.Printf("扫描广告活动数据失败: %v", err)
			continue
		}
		campaigns = append(campaigns, camp)
	}

	if err := rows.Err(); err != nil {
		log.Printf("遍历广告活动数据失败: %v", err)
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
			List:       campaigns,
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func GetCampaignByID(c *gin.Context) {
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
	var camp Campaign
	err = db.QueryRow(`
		SELECT id, name, description, status, budget, start_date, end_date 
		FROM ad_campaigns 
		WHERE id = ?
	`, id).Scan(&camp.ID, &camp.Name, &camp.Description, &camp.Status,
		&camp.Budget, &camp.StartDate, &camp.EndDate)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "活动不存在",
		})
		return
	}
	if err != nil {
		log.Printf("查询广告活动失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
		Data:    camp,
	})
}

func CreateCampaign(c *gin.Context) {
	var req CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "请求参数错误：名称、状态、预算和开始日期为必填项，状态只能为 active/paused/ended，预算不能为负数",
		})
		return
	}

	db := config.GetDB()
	result, err := db.Exec(`
		INSERT INTO ad_campaigns (name, description, status, budget, start_date, end_date)
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.Name, req.Description, req.Status, req.Budget, req.StartDate, req.EndDate)

	if err != nil {
		log.Printf("创建广告活动失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "创建广告活动失败，请稍后重试",
		})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "创建成功",
		Data: Campaign{
			ID:          id,
			Name:        req.Name,
			Description: req.Description,
			Status:      req.Status,
			Budget:      req.Budget,
			StartDate:   req.StartDate,
			EndDate:     req.EndDate,
		},
	})
}

func UpdateCampaign(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的 ID",
		})
		return
	}

	var req UpdateCampaignRequest
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

	if req.Name != nil {
		updateFields = append(updateFields, "name=?")
		args = append(args, *req.Name)
	}
	if req.Description != nil {
		updateFields = append(updateFields, "description=?")
		args = append(args, *req.Description)
	}
	if req.Status != nil {
		updateFields = append(updateFields, "status=?")
		args = append(args, *req.Status)
	}
	if req.Budget != nil {
		updateFields = append(updateFields, "budget=?")
		args = append(args, *req.Budget)
	}
	if req.StartDate != nil {
		updateFields = append(updateFields, "start_date=?")
		args = append(args, *req.StartDate)
	}
	if req.EndDate != nil {
		updateFields = append(updateFields, "end_date=?")
		args = append(args, *req.EndDate)
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "未提供任何更新字段",
		})
		return
	}

	// 检查状态值是否合法
	if req.Status != nil && *req.Status != "active" && *req.Status != "paused" && *req.Status != "ended" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "无效的状态值，可选值: active, paused, ended",
		})
		return
	}

	// 检查预算是否合法
	if req.Budget != nil && *req.Budget < 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "预算不能为负数",
		})
		return
	}

	sql := "UPDATE ad_campaigns SET " + strings.Join(updateFields, ", ") + " WHERE id=?"
	args = append(args, id)

	result, err := db.Exec(sql, args...)
	if err != nil {
		log.Printf("更新广告活动失败: %v", err)
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
			Message: "活动不存在",
		})
		return
	}

	// 查询更新后的完整记录
	var updatedCampaign Campaign
	err = db.QueryRow(`
		SELECT id, name, description, status, budget, start_date, end_date 
		FROM ad_campaigns 
		WHERE id = ?
	`, id).Scan(&updatedCampaign.ID, &updatedCampaign.Name, &updatedCampaign.Description,
		&updatedCampaign.Status, &updatedCampaign.Budget, &updatedCampaign.StartDate, &updatedCampaign.EndDate)

	if err != nil {
		log.Printf("查询更新后广告活动失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询更新结果失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "更新成功",
		Data:    updatedCampaign,
	})
}

func DeleteCampaign(c *gin.Context) {
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
	result, err := db.Exec("DELETE FROM ad_campaigns WHERE id = ?", id)
	if err != nil {
		log.Printf("删除广告活动失败: %v", err)
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
			Message: "活动不存在",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "删除成功",
	})
}