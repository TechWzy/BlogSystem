package rediss

import (
	g "Project/BlogSystem/internal/global"
	"Project/BlogSystem/internal/model"
	"Project/BlogSystem/internal/utils"
	"context"
	"fmt"
	"time"
)

// 由于本次项目的目的是练习业务逻辑开发的能力，所以不严格要求消息数据的存放位置，一律存放在redis数据库中

// InsertPrivateLetter 插入一条私信
func InsertPrivateLetter(letter *model.PrivateLetter) error {
	// 私信编号加一
	g.MessageID++
	letter.BaseMessage.CreatedTime = time.Now().Format("2006-01-02 15:04")
	// 设置键:"SRPrivate:sId:rId:id"
	key := fmt.Sprintf("%s#%v#%v#%v", "SRPrivate", letter.BaseMessage.SId, letter.BaseMessage.RId, letter.BaseMessage.ID)
	// 设置值:"subject:content"
	value := fmt.Sprintf("%s#%s#%s", letter.BaseMessage.CreatedTime, letter.Subject, letter.Content)
	t := time.Duration(letter.BaseMessage.Expiration)
	if err := g.Client.Set(context.Background(), key, value, t*time.Minute).Err(); err != nil {
		fmt.Println("无法向redis数据库插入SR类型私信...")
		return err
	}

	//fmt.Println("value:", value)

	// 插入键:"private:rId:sId:id",注意私信的过期日期暂时永久
	key = fmt.Sprintf("%s#%v#%v#%v", "RSPrivate", letter.BaseMessage.SId, letter.BaseMessage.RId, letter.BaseMessage.ID)
	if err := g.Client.Set(context.Background(), key, value, 0).Err(); err != nil {
		fmt.Println("无法向redis数据库插入RS类型私信...")
		return err
	}
	return nil // 插入成功
}

// GetPrivateLetterKeys 获取私信
func GetPrivateLetterKeys(letter *model.GetPrivateLetter) (keys []string, err error) {
	if letter.Mode == 1 {
		// 第一种查询模式:查询S对于R的所有私信
		patternString := utils.CombinePrivateLetterField(letter.SID, letter.RID)
		keys, err = g.Client.Keys(context.Background(), patternString).Result()
		return keys, err
	} else if letter.Mode == 2 {
		// 第二种模式查询:查询R的所有私信
		patternString := utils.CombineAllPrivateLetterField(letter.RID)
		keys, err = g.Client.Keys(context.Background(), patternString).Result()
		return keys, err
	} else {
		return []string{}, fmt.Errorf("fail to choose correct Mode")
	}
}

// GetPrivateLetterRecord 根据一个key值，获取一条私信的未解析字符串
func GetPrivateLetterRecord(key string) (v string, err error) {
	v, err = g.Client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	} else {
		return v, err
	}
}

// SetPrivateLetterExpiration 设置Private的过期期限
func SetPrivateLetterExpiration(key string, value string, exp int) {
	t := time.Duration(exp)
	g.Client.Set(context.Background(), key, value, t)
}
