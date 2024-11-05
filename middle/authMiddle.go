package middle

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"goUserCenter/config"
	"time"
)

// JWT 验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "未提供 token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.SecretKey), nil
		})

		if err != nil {
			c.JSON(401, gin.H{"error": "无效的token，解析错误：" + err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(401, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "无法获取token声明"})
			c.Abort()
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			c.JSON(401, gin.H{"error": "无法获取token过期时间"})
			c.Abort()
			return
		}

		if time.Now().Unix() > int64(exp) {
			c.JSON(401, gin.H{"error": "token已过期"})
			c.Abort()
			return
		}

		_, ok = claims["user_id"].(float64)
		if !ok {
			c.JSON(401, gin.H{"error": "无法获取用户ID"})
			c.Abort()
			return
		}

		c.Next()
	}
}
