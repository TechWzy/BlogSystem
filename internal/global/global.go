package global

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ArticleTagList []string

var (
	TagMap = make(map[string]bool)
	Db     *gorm.DB // Mysql数据库操作实例
	JwtKey = []byte("secret_key")
	Client = redis.NewClient(&redis.Options{
		Addr:     "192.168.246.128:6379", // Redis服务器地址
		Password: "",                     // Redis密码，如果没有密码则留空
		DB:       0,                      // 使用默认的数据库
		PoolSize: 20,
	})
	// TagLists 固定标签，为了方便文章管理，只有属于该标签的文章才能够被搜索到...
	TagLists = ArticleTagList{"大数据", "人工智能", "后端", "前端", "移动开发",
		"Java", "Python", "AIGC", "云平台", "C++", "C", "云原生", "工业互联网"}

	// MessageID 消息的ID编号
	MessageID int
)

const (
	Dsn = "root:123456@tcp(127.0.0.1:3306)/blogsystem?charset=utf8mb4&parseTime=True&loc=Local"
)
