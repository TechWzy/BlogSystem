package controller

import (
	"Project/BlogSystem/internal/dao/mysqls"
	"Project/BlogSystem/internal/dao/rediss"
	"Project/BlogSystem/internal/logic"
	"Project/BlogSystem/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
)

// CreateInfo 创建用户，自动分配账号
func CreateInfo(c *gin.Context) {
	var user model.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(200, gin.H{
			"msg": "Sex字只允许填man,woman,secret!",
		})
		c.Abort()
		return
	}

	c.Next()

	val, ok := c.Get("account")
	if ok {
		v := val.(string)
		user.Account = v
	}

	c.JSON(200, gin.H{
		"msg":  "success",
		"Info": user,
	})
	// 简单的业务逻辑，直接调用dao层
	mysqls.InsertInfo(&user)
	if err := rediss.RedisInsertData("AccountPassword", user.Account, user.Password); err != nil {
		fmt.Println("Fail to Insert a data into Redis!")
	}
}

// UserLogin 用户登录，将生成SignedToken
func UserLogin(c *gin.Context) {
	var loginInfo model.LoginInfo
	if err := c.ShouldBind(&loginInfo); err != nil {
		c.JSON(200, gin.H{
			"msg": "fail to bind!",
			"err": err,
		})
		c.Abort()
		return
	}
	password := logic.GetPassword(loginInfo.Account)
	if len(password) == 0 {
		c.JSON(200, gin.H{
			"msg": "账户不存在!",
		})
		c.Abort()
		return
	}
	if password != loginInfo.Password {
		c.JSON(200, gin.H{
			"msg": "密码不正确!",
		})
		c.Abort()
		return
	}

	// 当密码正确时，应该去获取该用户的一些基本信息
	authInfo := mysqls.GetAuthInfo(loginInfo.Account)

	// 生成Jwt
	c.Set("authInfo", authInfo)
	c.Next()
	v, ok := c.Get("TokenString")
	if !ok {
		fmt.Println("err:", "Fail to get the TokenString!")
		c.Abort()
		return
	}
	tokenString := v.(string)
	c.JSON(200, gin.H{
		"msg":         "登录成功!",
		"tokenString": tokenString,
	})
}
