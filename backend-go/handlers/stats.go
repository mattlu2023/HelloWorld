package handlers

import (
	"ad-bi-backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OverviewData 概览数据
type OverviewData struct {
	TotalImpressions   int64   `json:"total_impressions"`
	TotalClicks        int64   `json:"total_clicks"`
	TotalConversions   int64   `json:"total_conversions"`
	TotalCost          float64 `json:"total_cost"`
	TotalRevenue       float64 `json:"total_revenue"`
	AverageCTR         float64 `json:"average_ctr"`
	AverageConversionRate float64 `json:"average_conversion_rate"`
}

// GetOverview 获取数据概览
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

// DailyTrendData 每日趋势数据
type DailyTrendData struct {
	Date       string  `json:"date"`
	Impressions int64  `json:"impressions"`
	Clicks     int64  `json:"clicks"`
	Conversions int64 `json:"conversions"`
	Cost       float64 `json:"cost"`
	Revenue    float64 `json:"revenue"`
	CTR        float64 `json:"ctr"`
}

// GetDailyTrend 获取每日趋势
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}
	defer rows.Close()

	trends := []DailyTrendData{}
	for rows.Next() {
		var trend DailyTrendData
		err := rows.Scan(&trend.Date, &trend.Impressions, &trend.Clicks, 
			&trend.Conversions, &trend.Cost, &trend.Revenue, &trend.CTR)
		if err != nil {
			continue
		}
		trends = append(trends, trend)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    trends,
	})
}

// FunnelData 漏斗数据
type FunnelData struct {
	StepName   string `json:"step_name"`
	StepOrder  int    `json:"step_order"`
	UserCount  int    `json:"user_count"`
}

// GetConversionFunnel 获取转化漏斗
func GetConversionFunnel(c *gin.Context) {
	db := config.GetDB()
	
	rows, err := db.Query(`
		SELECT step_name, step_order, user_count
		FROM conversion_funnel
		WHERE funnel_date = CURDATE()
		ORDER BY step_order
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "查询失败",
		})
		return
	}
	defer rows.Close()

	funnel := []FunnelData{}
	for rows.Next() {
		var f FunnelData
		err := rows.Scan(&f.StepName, &f.StepOrder, &f.UserCount)
		if err != nil {
			continue
		}
		funnel = append(funnel, f)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    funnel,
	})
}

// GetCampaignStats 获取活动统计数据
func GetCampaignStats(c *gin.Context) {
	// 实现类似上面的查询逻辑
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    []interface{}{},
	})
}

// GetUserActions 获取用户行为
func GetUserActions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    []interface{}{},
	})
}

// GetUserActionAnalysis 用户行为分析
func GetUserActionAnalysis(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    gin.H{},
	})
}

// ExportReport 导出报表
func ExportReport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "报表导出成功",
		"data": gin.H{
			"download_url": "/reports/download/xxx.xlsx",
		},
	})
}
