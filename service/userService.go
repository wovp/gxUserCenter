package service

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"goUserCenter/config"
	"goUserCenter/gxmodule"
	"time"
)

// RegisterHandler 处理用户注册请求的函数
func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理用户注册逻辑
		username := c.PostForm("username")
		password := c.PostForm("password")
		email := c.PostForm("email")
		// 获取用户详细信息，这里使用默认值创建UserDetail实例
		userDetail, err := gxmodule.NewUserDetailRegister("", "", "", "2001-01-01", "", "", "", "", "")
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		newUser := gxmodule.NewUser(username, password, email)
		if err := newUser.Validate(); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := newUser.Save(db); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		// 设置用户详情信息中的UserID，使其与新创建的用户关联
		userDetail.UserID = newUser.ID
		// 保存用户详细信息到数据库
		if err := userDetail.Save(db); err != nil {
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

// getUserIdFromToken从token中获取用户ID（这里只是示例，你需要根据实际的token结构和解析方式来实现）
func getUserIdFromToken(c *gin.Context) int {
	// 假设token已经解析，并且claims中包含user_id字段
	claims, _ := getTokenClaims(c)
	return int(claims["user_id"].(float64))
}

// getTokenClaims获取token的声明内容（这里只是示例，你需要根据实际的token结构和解析方式来实现）
func getTokenClaims(c *gin.Context) (jwt.MapClaims, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("未提供token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("无效的token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("无法获取token声明")
	}

	return claims, nil
}

// parseDateOfBirth解析日期字符串为time.Time类型（假设日期格式为 "2000-01-02"）
func parseDateOfBirth(dateOfBirthStr string) time.Time {
	dateOfBirth, err := time.Parse("2006-01-02", dateOfBirthStr)
	if err != nil {
		panic(err)
	}

	return dateOfBirth
}

// GetUserDetailHandler获取用户详细信息的处理函数
func GetUserDetailHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从token或其他方式获取用户ID（这里假设从token中获取，你可以根据实际情况调整）
		userID := getUserIdFromToken(c)

		userDetail, err := gxmodule.GetUserDetailByUserID(db, userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"user_detail": userDetail})
	}
}

// UpdateUserDetailHandler更新用户详细信息的处理函数
func UpdateUserDetailHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从token获取用户ID
		userID := getUserIdFromToken(c)
		// 获取更新后的用户详细信息
		fullName := c.PostForm("full_name")
		phoneNumber := c.PostForm("phone_number")
		address := c.PostForm("address")
		dateOfBirthStr := c.PostForm("date_of_birth")
		gender := c.PostForm("gender")
		occupation := c.PostForm("occupation")
		avatar := c.PostForm("avatar")
		bio := c.PostForm("bio")
		school := c.PostForm("school")
		userDetail, err := gxmodule.GetUserDetailByUserID(db, userID)
		if err != nil {
			c.JSON(500, gin.H{"error_in_query": err.Error()})
			return
		}
		// 更新用户详细信息结构体的值
		userDetail.FullName = fullName
		userDetail.PhoneNumber = phoneNumber
		userDetail.Address = address
		if dateOfBirthStr == "" {
			dateOfBirthStr = "2001-01-01"
		}
		userDetail.DateOfBirth = parseDateOfBirth(dateOfBirthStr)
		userDetail.Gender = gender
		userDetail.Occupation = occupation
		userDetail.Avatar = avatar
		userDetail.Bio = bio
		userDetail.School = school
		err = userDetail.UpdateUserDetail(db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "用户详细信息更新成功"})
	}
}
