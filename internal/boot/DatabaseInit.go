package boot

import (
	"Project/BlogSystem/internal/dao/rediss"
	g "Project/BlogSystem/internal/global"
	"Project/BlogSystem/internal/model"
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitMysql 连接Mysql数据库
func InitMysql() {
	var err error
	g.Db, err = gorm.Open(mysql.Open(g.Dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	// 迁移数据库表结构
	err = g.Db.AutoMigrate(&model.User{})
	if err != nil {
		panic("Failed to migrate database table!")
	}

	err = g.Db.AutoMigrate(&model.Article{})
	if err != nil {
		panic("Failed to migrate database table!")
	}

	err = g.Db.AutoMigrate(&model.OtherArticleInfo{})
	if err != nil {
		panic("Failed to migrate database table!")
	}

	err = g.Db.AutoMigrate(&model.ArticleComment{})
	if err != nil {
		panic("Failed to migrate database table!")
	}

	err = g.Db.AutoMigrate(&model.InterestedTag{})
	if err != nil {
		panic("Failed to migrate database table!")
	}

	err = g.Db.AutoMigrate(&model.PassReaderRecord{})
	if err != nil {
		panic("Failed to migrate database table!")
	}
	fmt.Println("Mysql open and connect successfully!")
}

// InitRedis 初始化Redis实例化对象
func InitRedis() {
	// 连接Redis
	fmt.Println("Redis starting...")
	err := g.Client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis Open and connect successfully!")

	// 初始化一些数据类型
	if err := rediss.RedisInsertData("AccountPassword", "0", "0"); err != nil {
		fmt.Println("Fail to build set!")
	}

}

func CloseRedis() {
	err := g.Client.Close()
	if err != nil {
		fmt.Println("Fail to close connection pool:", err)
	}
}
