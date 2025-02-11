package util

import (
	"fmt"
	"mini-assets/internal/model"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// JWTClaims 定义了JWT的claims结构
type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.StandardClaims
}

// JWTSecret 用于签名JWT的密钥
var JWTSecret = []byte("your_secret_key") // 请替换为你的实际密钥

// GenerateToken 生成JWT令牌
func GenerateToken(user *model.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // 令牌过期时间
			Issuer:    "mini-assets",                         // 发行者
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken 验证JWT令牌
func ValidateToken(signedToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// AuthMiddleware 创建一个Gin中间件，用于验证JWT令牌
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// 确保Authorization字段是以"Bearer "开头的
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header must start with Bearer"})
			c.Abort()
			return
		}

		// 提取令牌字符串
		tokenString := authHeader[7:]

		// 解析令牌
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// 确保令牌使用的是HS256算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			// 将用户ID存储在上下文中，以便后续处理使用
			c.Set("userID", claims.UserID)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
		}
	}
}
