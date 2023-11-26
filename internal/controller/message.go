package controller

import (
	"Project/BlogSystem/internal/dao/mysqls"
	"Project/BlogSystem/internal/dao/rediss"
	g "Project/BlogSystem/internal/global"
	"Project/BlogSystem/internal/logic"
	"Project/BlogSystem/internal/model"
	"Project/BlogSystem/internal/service"
	"github.com/gin-gonic/gin"
)

// SendPrivateLetter 绑定私信内容
func SendPrivateLetter(c *gin.Context) {
	var privateLetter model.PrivateLetter
	if err := c.ShouldBind(&privateLetter); err != nil {
		c.JSON(200, gin.H{
			"Msg": "私信绑定失败...",
		})
		c.Abort()
		return
	}

	if flag := logic.IsUserExist(privateLetter.BaseMessage.RId); !flag {
		c.JSON(200, gin.H{
			"Msg": "Receiver的ID不存在...",
		})
		c.Abort()
		return
	}

	// 判断接收者是否存在
	if err := mysqls.JudgeUserIsExist(privateLetter.BaseMessage.RId); err != nil {
		c.JSON(200, gin.H{
			"Msg": "私信接收者不存在...",
		})
		c.Abort()
		return
	}

	v, ok := c.Get("authInfo")
	authInfo := v.(model.AuthInfo)
	if !ok {
		c.JSON(200, gin.H{
			"Msg": "无法解析authInfo",
		})
		c.Abort()
		return
	}

	// 补充数据
	privateLetter.BaseMessage.ID = g.MessageID
	privateLetter.BaseMessage.SId = authInfo.ID

	// 判断rId是否与rId相同
	if privateLetter.BaseMessage.SId == privateLetter.BaseMessage.RId {
		c.JSON(200, gin.H{
			"Msg": "私信发送者与接收者相同...",
		})
		c.Abort()
		return
	}

	// 插入redis数据库
	if err := rediss.InsertPrivateLetter(&privateLetter); err != nil {
		c.JSON(200, gin.H{
			"Msg": "PrivateLetter无法插入数据库...",
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"PrivateLetter": privateLetter,
	})
}

// GetPLetter 获取私信
func GetPLetter(c *gin.Context) {
	var privateLetterInfo model.GetPrivateLetter
	if err := c.ShouldBind(&privateLetterInfo); err != nil {
		c.JSON(200, gin.H{
			"Msg": "无法绑定信息...",
		})
		c.Abort()
		return
	}
	// 获取接收者信息
	v, ok := c.Get("authInfo")
	authInfo := v.(model.AuthInfo)
	if !ok {
		c.JSON(200, gin.H{
			"Msg": "无法解析authInfo...",
		})
		c.Abort()
		return
	}
	privateLetterInfo.RID = authInfo.ID

	// 判断senderId是否存在
	if err := mysqls.JudgeUserIsExist(privateLetterInfo.SID); err != nil {
		c.JSON(200, gin.H{
			"Msg": "Sender不存在...",
		})
		c.Abort()
		return
	}

	if authInfo.ID == privateLetterInfo.SID {
		c.JSON(200, gin.H{
			"Msg": "接收者与发送者相同...",
		})
		c.Abort()
		return
	}

	letterList, err := service.GetPrivateLetterList(&privateLetterInfo)
	if err != nil {
		c.JSON(200, gin.H{
			"Msg": "无法获取私信列表...",
		})
		c.Abort()
		return
	}
	var response model.PrivateLetterResponse
	if len(letterList) == 0 {
		response.Type = "空箱"
	} else if privateLetterInfo.Mode == 1 {
		response.Type = "来自" + letterList[0].Sender + "的私信"
		response.PrivateLetterList = letterList
	} else {
		response.Type = "用户的全部私信"
		response.PrivateLetterList = letterList
	}
	c.JSON(200, gin.H{
		"Response": response,
	})
}
