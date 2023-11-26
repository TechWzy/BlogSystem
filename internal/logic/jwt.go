package logic

import (
	g "Project/BlogSystem/internal/global"
	"Project/BlogSystem/internal/model"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

// GenTokenIA 根据业务需求封装一个JWT,此处仅仅包含ID与Account
func GenTokenIA(id int64, account string) string {
	// 创建自己的claims
	claims := model.MyClaims{
		ID:      id,
		Account: account,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "My-project",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 2)),
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的Secret签名并获得完整的编码后的字符串token
	signedToken, _ := token.SignedString(g.JwtKey)
	return signedToken
}

// AuthMiddlewareAI 令牌检验
func AuthMiddlewareAI() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// 验证令牌是否存在
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		// 从令牌中获取声明
		claims, err := ParseTokenIA(tokenString)
		if err != nil {
			fmt.Println("err:", err)
			c.Abort()
			return
		}

		// 将解析后的声明保存到上下文中，以便后续处理函数使用
		c.Set("claims", claims)
		c.Next()
	}
}

// ParseTokenIA 解析JWT
func ParseTokenIA(tokenString string) (*model.MyClaims, error) {
	// 使用ParseWithClaims方法解析结构体
	token, err := jwt.ParseWithClaims(tokenString, &model.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return g.JwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	// 对token对象中的claims进行断言
	if claims, ok := token.Claims.(*model.MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// GetAuthTokenMiddleware 生成登录成功后的tokenString
func GetAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("authInfo")
		if !ok {
			c.Set("TokenString", "")
			c.Abort()
			return
		}
		authInfo := v.(model.AuthInfo)
		// 生成claims
		claims := model.AuthInfoClaims{
			AuthInfo: authInfo,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "Login-Jwt",
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(g.JwtKey)
		c.Set("TokenString", tokenString)
	}
}

// ParseAuthInfoMiddleware 检验登录成功后的令牌
func ParseAuthInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// 验证令牌是否存在
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			fmt.Println("tokenString Obtainment fail")
			c.Abort()
			return
		}

		// 解析和验证令牌
		token, err := jwt.ParseWithClaims(tokenString, &model.AuthInfoClaims{}, func(token *jwt.Token) (interface{}, error) {
			return g.JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		// 从令牌中提取声明
		claims, ok := token.Claims.(*model.AuthInfoClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token claims",
			})
			c.Abort()
			return
		}
		// 将解析后的声明保存到上下文中，以便后续处理函数使用
		c.Set("authInfo", claims.AuthInfo)
		c.Next()
	}
}
