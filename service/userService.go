package service

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"goUserCenter/gxmodule"
)

// RegisterHandler 处理用户注册请求的函数
func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理用户注册逻辑
		username := c.PostForm("username")
		password := c.PostForm("password")
		email := c.PostForm("email")

		newUser := gxmodule.NewUser(username, password, email)
		if err := newUser.Validate(); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := newUser.Save(db); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{"message": "用户创建成功"})
	}
}

// LoginHandler 处理用户登录请求的函数
func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理用户登录逻辑
		username := c.PostForm("username")
		password := c.PostForm("password")

		user := gxmodule.User{Username: username, Password: password}
		token, authenticated := user.Authenticate(db)
		if authenticated {
			c.JSON(200, gin.H{"token": token})
		} else {
			c.JSON(401, gin.H{"error": "登录失败"})
		}
	}
}
