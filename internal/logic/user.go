package logic

import (
	"Project/BlogSystem/internal/dao/mysqls"
	"Project/BlogSystem/internal/dao/rediss"
	"Project/BlogSystem/internal/utils"
)

// 获取独立的账号account
func getAccount(set string) string {
	account := utils.GenerateRandomString(8)
	_, err := rediss.GetAccountPassword(set, account)
	if err == nil {
		return getAccount(set)
	}
	return account
}

// GetPassword 输入账号获取密码
func GetPassword(account string) string {
	password, _ := rediss.GetAccountPassword("AccountPassword", account)
	return password
}

func IsUserExist(userId int) bool {
	if err := mysqls.JudgeUserIsExist(userId); err != nil {
		return false
	} else {
		return true
	}
}
