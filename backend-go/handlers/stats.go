package handlers

import (
	"ad-bi-backend/config"
	"bytes"
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type OverviewData struct {
	TotalImpressions      int64   `json:"total_impressions"`
	TotalClicks           int64   `json:"total_clicks"`
	TotalConversions      int64   `json:"total_conversions"`
	TotalCost             float64 `json:"total_cost"`
	TotalRevenue          float64 `json:"total_revenue"`
	AverageCTR            float64 `json:"average_ctr"`
	AverageConversionRate float64 `json:"average_conversion_rate"`
}

func GetOverview(c *gin.Context) {
	db := config.GetDB()

	var data OverviewData
	err := db.QueryRow(`
		SELECT 
			COALESCE(SUM(impressions), 0),
			COALESCE(SUM(clicks), 0),
			COALESCE(SUM(conversions), 0),
			COALESCE(SUM(cost), 0),
			COALESCE(SUM(revenue), 0),
			COALESCE(SUM(clicks) * 100.0 / NULLIF(SUM(impressions), 0), 0),
			COALESCE(SUM(conversions) * 100.0 / NULLIF(SUM(clicks), 0), 0)
		FROM ad_stats_daily
	`).Scan(&data.TotalImpressions, &data.TotalClicks, &data.TotalConversions,
		&data.TotalCost, &data.TotalRevenue, &data.AverageCTR, &data.AverageConversionRate)

	if err != nil {
		log.Printf("查询数据概览失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

type DailyTrendData struct {
	Date        string  `json:"date"`
	Impressions int64   `json:"impressions"`
	Clicks      int64   `json:"clicks"`
	Conversions int64   `json:"conversions"`
	Cost        float64 `json:"cost"`
	Revenue     float64 `json:"revenue"`
	CTR         float64 `json:"ctr"`
}

func GetDailyTrend(c *gin.Context) {
	db := config.GetDB()

	rows, err := db.Query(`
		SELECT 
			stat_date,
			SUM(impressions),
			SUM(clicks),
			SUM(conversions),
			SUM(cost),
			SUM(revenue),
			COALESCE(SUM(clicks) * 100.0 / NULLIF(SUM(impressions), 0), 0)
		FROM ad_stats_daily
		GROUP BY stat_date
		ORDER BY stat_date DESC
		LIMIT 30
	`)
	if err != nil {
		log.Printf("查询每日趋势失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}
	defer rows.Close()

	trends := []DailyTrendData{}
	for rows.Next() {
		var trend DailyTrendData
		if err := rows.Scan(&trend.Date, &trend.Impressions, &trend.Clicks,
			&trend.Conversions, &trend.Cost, &trend.Revenue, &trend.CTR); err != nil {
			log.Printf("扫描每日趋势数据失败: %v", err)
			continue
		}
		trends = append(trends, trend)
	}

	if err := rows.Err(); err != nil {
		log.Printf("遍历每日趋势数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
		Data:    trends,
	})
}

type FunnelData struct {
	StepName  string `json:"step_name"`
	StepOrder int    `json:"step_order"`
	UserCount int    `json:"user_count"`
}

func GetConversionFunnel(c *gin.Context) {
	db := config.GetDB()

	rows, err := db.Query(`
		SELECT step_name, step_order, user_count
		FROM conversion_funnel
		WHERE funnel_date = CURDATE()
		ORDER BY step_order
	`)
	if err != nil {
		log.Printf("查询转化漏斗失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}
	defer rows.Close()

	funnel := []FunnelData{}
	for rows.Next() {
		var f FunnelData
		if err := rows.Scan(&f.StepName, &f.StepOrder, &f.UserCount); err != nil {
			log.Printf("扫描转化漏斗数据失败: %v", err)
			continue
		}
		funnel = append(funnel, f)
	}

	if err := rows.Err(); err != nil {
		log.Printf("遍历转化漏斗数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
		Data:    funnel,
	})
}

type CampaignStatsData struct {
	Date        string  `json:"date"`
	Impressions int64   `json:"impressions"`
	Clicks      int64   `json:"clicks"`
	Conversions int64   `json:"conversions"`
	Cost        float64 `json:"cost"`
	Revenue     float64 `json:"revenue"`
}

func GetCampaignStats(c *gin.Context) {
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

	rows, err := db.Query(`
		SELECT 
			stat_date,
			SUM(impressions),
			SUM(clicks),
			SUM(conversions),
			SUM(cost),
			SUM(revenue)
		FROM ad_stats_daily
		WHERE campaign_id = ?
		GROUP BY stat_date
		ORDER BY stat_date DESC
		LIMIT 30
	`, id)
	if err != nil {
		log.Printf("查询活动统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}
	defer rows.Close()

	stats := []CampaignStatsData{}
	for rows.Next() {
		var stat CampaignStatsData
		if err := rows.Scan(&stat.Date, &stat.Impressions, &stat.Clicks,
			&stat.Conversions, &stat.Cost, &stat.Revenue); err != nil {
			log.Printf("扫描活动统计数据失败: %v", err)
			continue
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		log.Printf("遍历活动统计数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
		Data:    stats,
	})
}

type UserActionData struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Action    string `json:"action"`
	Page      string `json:"page"`
	CreatedAt string `json:"created_at"`
}

func GetUserActions(c *gin.Context) {
	db := config.GetDB()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 查询总数
	var total int64
	err := db.QueryRow("SELECT COUNT(*) FROM user_actions").Scan(&total)
	if err != nil {
		log.Printf("查询用户行为总数失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}

	// 计算总页数
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	rows, err := db.Query(`
		SELECT id, user_id, action, page, created_at
		FROM user_actions
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, pageSize, offset)
	if err != nil {
		log.Printf("查询用户行为失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}
	defer rows.Close()

	actions := []UserActionData{}
	for rows.Next() {
		var action UserActionData
		if err := rows.Scan(&action.ID, &action.UserID, &action.Action,
			&action.Page, &action.CreatedAt); err != nil {
			log.Printf("扫描用户行为数据失败: %v", err)
			continue
		}
		actions = append(actions, action)
	}

	if err := rows.Err(); err != nil {
		log.Printf("遍历用户行为数据失败: %v", err)
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
			List:       actions,
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

type UserActionAnalysisData struct {
	Action     string `json:"action"`
	Count      int64  `json:"count"`
	Percentage float64 `json:"percentage"`
}

func GetUserActionAnalysis(c *gin.Context) {
	db := config.GetDB()

	rows, err := db.Query(`
		SELECT 
			action,
			COUNT(*) as count
		FROM user_actions
		GROUP BY action
		ORDER BY count DESC
	`)
	if err != nil {
		log.Printf("查询用户行为分析失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}
	defer rows.Close()

	var total int64 = 0
	analysis := []UserActionAnalysisData{}
	for rows.Next() {
		var data UserActionAnalysisData
		if err := rows.Scan(&data.Action, &data.Count); err != nil {
			log.Printf("扫描用户行为分析数据失败: %v", err)
			continue
		}
		total += data.Count
		analysis = append(analysis, data)
	}

	if err := rows.Err(); err != nil {
		log.Printf("遍历用户行为分析数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询失败",
		})
		return
	}

	for i := range analysis {
		if total > 0 {
			analysis[i].Percentage = float64(analysis[i].Count) * 100.0 / float64(total)
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    0,
		Message: "success",
		Data:    analysis,
	})
}

func ExportReport(c *gin.Context) {
	db := config.GetDB()

	// 获取日期范围参数
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	// 构建查询条件
	var query string
	var args []interface{}

	if startDate != "" && endDate != "" {
		query = `
			SELECT 
				stat_date,
				campaign_id,
				SUM(impressions) as impressions,
				SUM(clicks) as clicks,
				SUM(conversions) as conversions,
				SUM(cost) as cost,
				SUM(revenue) as revenue,
				COALESCE(SUM(clicks) * 100.0 / NULLIF(SUM(impressions), 0), 0) as ctr,
				COALESCE(SUM(conversions) * 100.0 / NULLIF(SUM(clicks), 0), 0) as conversion_rate
			FROM ad_stats_daily
			WHERE stat_date BETWEEN ? AND ?
			GROUP BY stat_date, campaign_id
			ORDER BY stat_date DESC, campaign_id
		`
		args = append(args, startDate, endDate)
	} else if startDate != "" {
		query = `
			SELECT 
				stat_date,
				campaign_id,
				SUM(impressions) as impressions,
				SUM(clicks) as clicks,
				SUM(conversions) as conversions,
				SUM(cost) as cost,
				SUM(revenue) as revenue,
				COALESCE(SUM(clicks) * 100.0 / NULLIF(SUM(impressions), 0), 0) as ctr,
				COALESCE(SUM(conversions) * 100.0 / NULLIF(SUM(clicks), 0), 0) as conversion_rate
			FROM ad_stats_daily
			WHERE stat_date >= ?
			GROUP BY stat_date, campaign_id
			ORDER BY stat_date DESC, campaign_id
		`
		args = append(args, startDate)
	} else {
		query = `
			SELECT 
				stat_date,
				campaign_id,
				SUM(impressions) as impressions,
				SUM(clicks) as clicks,
				SUM(conversions) as conversions,
				SUM(cost) as cost,
				SUM(revenue) as revenue,
				COALESCE(SUM(clicks) * 100.0 / NULLIF(SUM(impressions), 0), 0) as ctr,
				COALESCE(SUM(conversions) * 100.0 / NULLIF(SUM(clicks), 0), 0) as conversion_rate
			FROM ad_stats_daily
			GROUP BY stat_date, campaign_id
			ORDER BY stat_date DESC, campaign_id
		`
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("查询报表数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "查询报表数据失败",
		})
		return
	}
	defer rows.Close()

	// 创建CSV缓冲区
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// 写入CSV表头
	headers := []string{
		"日期", "活动ID", "展示量", "点击量", "转化量", "成本", "收入", "点击率(%)", "转化率(%)",
	}
	if err := writer.Write(headers); err != nil {
		log.Printf("写入CSV表头失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "生成CSV失败",
		})
		return
	}

	// 写入数据行
	for rows.Next() {
		var (
			statDate       string
			campaignID     int64
			impressions    int64
			clicks         int64
			conversions    int64
			cost           float64
			revenue        float64
			ctr            float64
			conversionRate float64
		)
		if err := rows.Scan(&statDate, &campaignID, &impressions, &clicks,
			&conversions, &cost, &revenue, &ctr, &conversionRate); err != nil {
			log.Printf("扫描报表数据失败: %v", err)
			continue
		}

		record := []string{
			statDate,
			strconv.FormatInt(campaignID, 10),
			strconv.FormatInt(impressions, 10),
			strconv.FormatInt(clicks, 10),
			strconv.FormatInt(conversions, 10),
			strconv.FormatFloat(cost, 'f', 2, 64),
			strconv.FormatFloat(revenue, 'f', 2, 64),
			strconv.FormatFloat(ctr, 'f', 2, 64),
			strconv.FormatFloat(conversionRate, 'f', 2, 64),
		}

		if err := writer.Write(record); err != nil {
			log.Printf("写入CSV记录失败: %v", err)
			c.JSON(http.StatusInternalServerError, APIResponse{
				Code:    500,
				Message: "生成CSV失败",
			})
			return
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("遍历报表数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "遍历报表数据失败",
		})
		return
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Printf("刷新CSV写入器失败: %v", err)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "生成CSV失败",
		})
		return
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := "report_" + timestamp + ".csv"

	// 设置响应头
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Length", strconv.Itoa(buf.Len()))

	// 返回CSV数据
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}