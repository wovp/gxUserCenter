package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"goUserCenter/middle"
	"goUserCenter/service"
)

// SetupRoutes 用于设置应用的所有路由
func SetupRoutes(db *sql.DB) *gin.Engine {
	router := gin.Default()

	// 用户注册路由
	router.POST("/register", service.RegisterHandler(db))

	// 用户登录路由
	router.POST("/login", service.LoginHandler(db))

	// 需要认证的路由，添加JWT验证中间件
	authenticatedRoutes := router.Group("/user")
	authenticatedRoutes.Use(middle.AuthMiddleware())

	// 示例：返回 "Hello, World!" 的路由，需登录后访问
	authenticatedRoutes.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})
	userQuery := authenticatedRoutes.Group("/query")
	userUpdate := authenticatedRoutes.Group("/update")

	// 获取用户详细信息路由
	userQuery.GET("/user_detail", service.GetUserDetailHandler(db))

	// 更新用户详细信息路由
	userUpdate.POST("/user_detail", service.UpdateUserDetailHandler(db))

	return router
}
