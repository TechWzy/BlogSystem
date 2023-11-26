package controller

import (
	"Project/BlogSystem/internal/dao/mysqls"
	"Project/BlogSystem/internal/dao/rediss"
	"Project/BlogSystem/internal/logic"
	"Project/BlogSystem/internal/model"
	"Project/BlogSystem/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

// SetFollowCond 设置粉丝关注状态
func SetFollowCond(c *gin.Context) {
	var fCond model.SetFollowCondition
	if err := c.ShouldBind(&fCond); err != nil {
		fmt.Println("BindingError:", err.Error())
		c.JSON(200, gin.H{
			"Msg": "绑定失败...",
		})
		c.Abort()
		return
	}

	// 判断idolId的合法性
	if !logic.IsUserExist(fCond.IdolId) {
		c.JSON(200, gin.H{
			"Msg": "idolId不存在...",
		})
		c.Abort()
		return
	}

	v, _ := c.Get("authInfo")
	authInfo := v.(model.AuthInfo)
	fCond.FollowId = authInfo.ID

	// 判断idolId是否为本人
	if fCond.FollowId == fCond.IdolId {
		c.JSON(200, gin.H{
			"Msg": "FollowId与IdolId相同...",
		})
		c.Abort()
		return
	}

	// 根据cond划分为关注与取消关注两种业务
	if fCond.Cond {
		res, e := service.AddFollow(fCond.IdolId, fCond.FollowId)
		if e != nil {
			c.JSON(200, gin.H{
				"Msg": res,
			})
			c.Abort()
			return
		} else {
			amount, _ := rediss.GetFollowAmount(fCond.IdolId)
			name, _ := mysqls.GetSenderNameById(fCond.IdolId)
			res := "Idol" + name + "当前的粉丝数量为:" + amount
			c.JSON(200, gin.H{
				"Msg": "添加成功...",
				"Res": res,
			})
		}
	} else {
		res, e := service.DelFollow(fCond.IdolId, fCond.FollowId)
		if e != nil {
			c.JSON(200, gin.H{
				"Msg": res,
			})
			c.Abort()
			return
		} else {
			amount, _ := rediss.GetFollowAmount(fCond.IdolId)
			name, _ := mysqls.GetSenderNameById(fCond.IdolId)
			res := "Idol" + name + "当前的粉丝数量为:" + amount
			c.JSON(200, gin.H{
				"Msg": "删除成功...",
				"Res": res,
			})
		}
	}
}
