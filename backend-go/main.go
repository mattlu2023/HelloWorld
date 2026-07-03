package main

// 广告BI系统后端服务入口文件，基于Gin框架构建RESTful API

import (
	"ad-bi-backend/config"
	"ad-bi-backend/handlers"
	"ad-bi-backend/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 入口阶段只做固定成本的初始化工作：加载配置、建立数据库连接、注册路由。
	// 这里不承载业务逻辑，便于后续定位性能瓶颈时把注意力集中在 handler 和数据库查询上。
	cfg := config.LoadConfig()

	// 启动时预先验证数据库可用性，避免服务监听成功后才在首个请求上暴露连接问题。
	if err := config.InitDatabase(cfg); err != nil {
		log.Fatal("数据库初始化失败:", err)
	}

	// 使用 Gin 默认中间件栈，统一提供日志和 panic 恢复能力。
	r := gin.Default()

	// CORS 作为全局中间件挂载，保证浏览器请求在进入业务处理前完成跨域校验。
	r.Use(middleware.CORSMiddleware())

	// 健康检查保持匿名可访问，供容器探针和外部监控直接探活。
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "广告 BI 系统运行正常",
		})
	})

	// 所有业务接口统一挂在 /api/v1 下，便于版本演进和网关路由管理。
	api := r.Group("/api/v1")
	{
		// 登录和注册必须保持匿名可访问，避免被认证中间件拦截。
		api.POST("/login", handlers.Login)
		api.POST("/register", handlers.Register)

		// 其余业务接口默认要求认证；新的受保护接口应优先注册到该分组下。
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			// 广告活动管理。
			auth.GET("/campaigns", handlers.GetCampaigns)
			auth.GET("/campaigns/:id", handlers.GetCampaignByID)
			auth.POST("/campaigns", handlers.CreateCampaign)
			auth.PUT("/campaigns/:id", handlers.UpdateCampaign)
			auth.DELETE("/campaigns/:id", handlers.DeleteCampaign)

			// 广告单元管理。
			auth.GET("/ad-units", handlers.GetAdUnits)
			auth.GET("/ad-units/:id", handlers.GetAdUnitByID)
			auth.POST("/ad-units", handlers.CreateAdUnit)
			auth.PUT("/ad-units/:id", handlers.UpdateAdUnit)
			auth.DELETE("/ad-units/:id", handlers.DeleteAdUnit)

			// 数据统计接口。
			auth.GET("/stats/overview", handlers.GetOverview)
			auth.GET("/stats/daily-trend", handlers.GetDailyTrend)
			auth.GET("/stats/campaign/:id", handlers.GetCampaignStats)
			auth.GET("/stats/funnel", handlers.GetConversionFunnel)

			// 用户行为分析。
			auth.GET("/user-actions", handlers.GetUserActions)
			auth.GET("/user-actions/analysis", handlers.GetUserActionAnalysis)

			// 报表导出。
			auth.GET("/reports/export", handlers.ExportReport)
		}
	}

	// 监听端口后进入请求处理阶段；正常运行时主 goroutine 会阻塞在这里。
	log.Printf("服务器启动在端口 %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
