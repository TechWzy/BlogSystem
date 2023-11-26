package rediss

import (
	g "Project/BlogSystem/internal/global"
	"context"
	"fmt"
)

// IsFollow 判断a是否为b的粉丝,在添加或删除粉丝数据前
func IsFollow(idolId int, followId int) (bool, error) {
	key := "Idol"
	key = fmt.Sprintf("%s:%v", key, idolId)

	// 调用 SIsMember命令，判断某一个元素是否属于该集合(假定数据库不会出现连接超时的情况)
	result, err := g.Client.SIsMember(context.Background(), key, followId).Result()

	// 如果集合不存在，将返回redis.Nil
	if err != nil {
		return false, err
	} else {
		return result, err
	}
}

// AddFollow 添加粉丝
func AddFollow(idolId int, followId int) error {
	key := "Idol"
	key = fmt.Sprintf("%s:%v", key, idolId)

	if err := g.Client.SAdd(context.Background(), key, followId).Err(); err != nil {
		return err
	} else {
		return nil
	}
}

// DeleteFollow 删除粉丝
func DeleteFollow(idolId int, followId int) error {
	key := "Idol"
	key = fmt.Sprintf("%s:%v", key, idolId)
	err := g.Client.SRem(context.Background(), key, followId).Err()
	return err
}

// IsIdolKeyExist 判断NumberIdol:id是否存在
func IsIdolKeyExist(idolId int) error {
	key := "NumberIdol"
	key = fmt.Sprintf("%s:%v", key, idolId)

	// 判断键是否存在
	exist, err := g.Client.Exists(context.Background(), key).Result()

	// 对err进行处理
	if err != nil {
		return err
	}

	if exist == 1 {
		return nil
	} else {
		if e := g.Client.Set(context.Background(), key, 0, 0).Err(); e != nil {
			return e
		}
	}
	return nil
}

// AdjustFollowNumber 调整粉丝数量
func AdjustFollowNumber(idolId int, increment int) error {
	key := "NumberIdol"
	key = fmt.Sprintf("%s:%v", key, idolId)

	if increment == 1 {
		e := g.Client.Incr(context.Background(), key).Err()
		if e != nil {
			return e
		}
	} else if increment == -1 {
		e := g.Client.Decr(context.Background(), key).Err()
		if e != nil {
			return e
		}
	} else {
		return fmt.Errorf("wrong increment")
	}
	return nil
}

// GetFollowAmount 获取粉丝数量
func GetFollowAmount(idolId int) (string, error) {
	key := "NumberIdol"
	key = fmt.Sprintf("%s:%v", key, idolId)
	v, e := g.Client.Get(context.Background(), key).Result()
	if e != nil {
		return "无法获取粉丝数量...", e
	} else {
		return v, e
	}
}

// GetFollowList 获取粉丝列表
func GetFollowList(idolId int) ([]string, error) {
	key := "Idol"
	key = fmt.Sprintf("%s:%v", key, idolId)

	// 获取列表,SMembers只返回字符串类型
	result, err := g.Client.SMembers(context.Background(), key).Result()
	if err != nil {
		return []string{}, err
	}
	return result, err
}
