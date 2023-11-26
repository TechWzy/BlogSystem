package rediss

import (
	g "Project/BlogSystem/internal/global"
	"Project/BlogSystem/internal/utils"
	"context"
	"fmt"
)

// RedisInsertData 插入account:password,set代表集合的名称
func RedisInsertData(set string, account string, password string) error {
	// 将数据格式化为account:password字符串
	err := g.Client.SAdd(context.Background(), set, fmt.Sprintf("%s:%s", account, password)).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetAccountPassword 查询账户密码
func GetAccountPassword(set string, account string) (string, error) {
	values, err := g.Client.SMembers(context.Background(), set).Result()
	if err != nil {
		return "", err
	}
	for _, value := range values {
		keyVal := utils.ParseKeyValue(value)
		if keyVal.Key == account {
			return keyVal.Value, nil
		}
	}
	return "", fmt.Errorf(" account not found")
}
