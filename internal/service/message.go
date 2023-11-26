package service

import (
	"Project/BlogSystem/internal/dao/mysqls"
	"Project/BlogSystem/internal/dao/rediss"
	"Project/BlogSystem/internal/model"
	"strconv"
	"strings"
)

// 获取S向R发送的所有私信

func GetPrivateLetterList(letterInfo *model.GetPrivateLetter) ([]model.PrivateLetterField, error) {
	keys, err := rediss.GetPrivateLetterKeys(letterInfo)
	if err != nil {
		return []model.PrivateLetterField{}, err
	}
	// 已知keys列表里面的元素，能够获得该keys指定的私信
	var letterList []model.PrivateLetterField
	for _, key := range keys {
		// key代表每一条私信的键值
		v, e := rediss.GetPrivateLetterRecord(key)
		if e != nil {
			continue
		}
		// 分解key值,获取发送者信息
		parts := strings.Split(key, "#")
		senderId, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		var sender string
		sender, err = mysqls.GetSenderNameById(senderId)
		parts = strings.Split(v, "#")
		var letter model.PrivateLetterField
		letter.CreatedTime = parts[0]
		letter.Subject, _ = strconv.Unquote(`"` + parts[1] + `"`) // linux编码转化中文
		letter.Content, _ = strconv.Unquote(`"` + parts[2] + `"`)
		letter.SenderId = senderId
		letter.Sender = sender
		letterList = append(letterList, letter)
		if letterInfo.Mode == 2 {
			continue
		} else if letterInfo.Mode == 1 && letterInfo.Exp == 0 {
			continue
		} else {
			rediss.SetPrivateLetterExpiration(key, v, letterInfo.Exp)
		}
	}
	return letterList, nil
}
