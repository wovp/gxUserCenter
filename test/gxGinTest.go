package test

import (
	"github.com/gin-gonic/gin"
	user "goUserCenter/gxmodule"
	"goUserCenter/middle"
	"log"
)

func G() {
	// 连接数据库
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	router.POST("/register", func(c *gin.Context) {
		// 处理用户注册逻辑
		username := c.PostForm("username")
		password := c.PostForm("password")
		email := c.PostForm("email")

		newUser := user.NewUser(username, password, email)
		if err := newUser.Validate(); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := newUser.Save(db); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{"message": "用户创建成功"})
	})

	router.POST("/login", func(c *gin.Context) {
		// 处理用户登录逻辑
		username := c.PostForm("username")
		password := c.PostForm("password")

		user1 := user.User{Username: username, Password: password}
		token, authenticated := user1.Authenticate(db)
		if authenticated {
			c.JSON(200, gin.H{"token": token})
		} else {
			c.JSON(401, gin.H{"error": "登录失败"})
		}
	})

	router.GET("/hello", middle.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	err = router.Run(":8080")
	if err != nil {
		return
	}
}
