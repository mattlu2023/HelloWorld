package main

import (
	"ad-bi-backend/config"
	"ad-bi-backend/handlers"
	"ad-bi-backend/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	if err := config.InitDatabase(cfg); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 配置 CORS
	r.Use(middleware.CORSMiddleware())

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "广告 BI 系统运行正常",
		})
	})

	// API 路由组
	api := r.Group("/api/v1")
	{
		// 用户相关接口
		api.POST("/login", handlers.Login)
		api.POST("/register", handlers.Register)

		// 需要认证的接口
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			// 广告活动管理
			auth.GET("/campaigns", handlers.GetCampaigns)
			auth.GET("/campaigns/:id", handlers.GetCampaignByID)
			auth.POST("/campaigns", handlers.CreateCampaign)
			auth.PUT("/campaigns/:id", handlers.UpdateCampaign)
			auth.DELETE("/campaigns/:id", handlers.DeleteCampaign)

			// 广告单元管理
			auth.GET("/ad-units", handlers.GetAdUnits)
			auth.GET("/ad-units/:id", handlers.GetAdUnitByID)
			auth.POST("/ad-units", handlers.CreateAdUnit)
			auth.PUT("/ad-units/:id", handlers.UpdateAdUnit)
			auth.DELETE("/ad-units/:id", handlers.DeleteAdUnit)

			// 数据统计接口
			auth.GET("/stats/overview", handlers.GetOverview)
			auth.GET("/stats/daily-trend", handlers.GetDailyTrend)
			auth.GET("/stats/campaign/:id", handlers.GetCampaignStats)
			auth.GET("/stats/funnel", handlers.GetConversionFunnel)

			// 用户行为分析
			auth.GET("/user-actions", handlers.GetUserActions)
			auth.GET("/user-actions/analysis", handlers.GetUserActionAnalysis)

			// 报表导出
			auth.GET("/reports/export", handlers.ExportReport)
		}
	}

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
