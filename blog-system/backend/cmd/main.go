package main

import (
	"log"
	"os"

	"blog-system/internal/config"
	"blog-system/internal/handler"
	"blog-system/internal/middleware"
	"blog-system/internal/repository"
	"blog-system/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	// 初始化数据库
	if err := repository.InitDB(&cfg.Database); err != nil {
		log.Fatal("初始化数据库失败:", err)
	}

	// 初始化 Redis
	if err := repository.InitRedis(&cfg.Redis); err != nil {
		log.Fatal("初始化 Redis 失败:", err)
	}

	// 初始化短信服务
	smsService, err := service.NewSMSService(&cfg.SMS)
	if err != nil {
		log.Fatal("初始化短信服务失败:", err)
	}

	// 创建 Gin 引擎
	r := gin.Default()

	// 使用中间件
	r.Use(middleware.CORS())

	// 静态文件
	r.Static("/static", "./static")
	r.Static("/uploads/images", "./uploads/images")
	r.Static("/uploads/pdf", "./uploads/pdf")
	r.Static("/uploads/video", "./uploads/video")

	// 初始化处理器
	authHandler := handler.NewAuthHandler()
	authHandler.SetSMSService(smsService)
	postHandler := handler.NewPostHandler()
	commentHandler := handler.NewCommentHandler()
	uploadHandler := handler.NewUploadHandler()
	adminHandler := handler.NewAdminHandler()

	// API 路由组
	api := r.Group("/api")
	{
		// 认证相关
		auth := api.Group("/auth")
		{
			auth.POST("/send-sms", authHandler.SendSMSCode)
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.RegisterAdmin)
		}

		// 公开接口
		public := api.Group("")
		{
			// 博文
			public.GET("/posts", postHandler.GetPostList)
			public.GET("/posts/:slug", postHandler.GetPost)
			public.GET("/posts/search", postHandler.SearchPosts)
			public.GET("/posts/hot", postHandler.GetHotPosts)
			public.POST("/posts/:id/like", postHandler.LikePost)
			public.POST("/posts/:id/share", postHandler.SharePost)

			// 评论
			public.GET("/posts/:post_id/comments", commentHandler.GetComments)
			public.POST("/comments", commentHandler.CreateComment)
			public.POST("/comments/:id/like", commentHandler.LikeComment)

			// 上传
			public.POST("/upload/image", uploadHandler.UploadImage)
			public.POST("/upload/pdf", uploadHandler.UploadPDF)
			public.POST("/upload/video", uploadHandler.UploadVideo)
		}

		// 需要认证的接口
		protected := api.Group("")
		protected.Use(middleware.JWTAuth())
		{
			// 用户信息
			protected.GET("/user/info", authHandler.GetUserInfo)
			protected.POST("/user/logout", authHandler.Logout)

			// 博文管理
			protected.POST("/posts", postHandler.CreatePost)
			protected.PUT("/posts/:id", postHandler.UpdatePost)
			protected.DELETE("/posts/:id", postHandler.DeletePost)

			// 评论管理
			protected.GET("/admin/comments/pending", commentHandler.GetPendingComments)
			protected.POST("/admin/comments/audit", commentHandler.AuditComment)
			protected.DELETE("/comments/:id", commentHandler.DeleteComment)

			// 上传管理
			protected.GET("/uploads", uploadHandler.GetUploadList)
			protected.DELETE("/uploads/:id", uploadHandler.DeleteUpload)

			// 后台管理
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminOnly())
			{
				admin.GET("/dashboard", adminHandler.GetDashboard)
				
				// 敏感词管理
				admin.GET("/sensitive", adminHandler.GetSensitiveList)
				admin.POST("/sensitive", adminHandler.AddSensitiveWord)
				admin.DELETE("/sensitive/:id", adminHandler.DeleteSensitiveWord)
				admin.PUT("/sensitive/status", adminHandler.UpdateSensitiveStatus)
				admin.POST("/sensitive/import", adminHandler.ImportSensitiveWords)

				// 统计
				admin.GET("/statistics", adminHandler.GetStatistics)

				// 备份
				admin.POST("/backup", adminHandler.BackupDatabase)
			}
		}
	}

	// 启动服务器
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在端口 %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}
