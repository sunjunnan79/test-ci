package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}

type Config struct {
	MySQL struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"mysql"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
	} `yaml:"redis"`
	Kafka struct {
		Brokers []string `yaml:"brokers"`
	} `yaml:"kafka"`
}

func loadConfig(filePath string) (*Config, error) {
	// 打开配置文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// 解析配置文件
	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	// 启动 Gin 服务
	r := setupRouter(InitConfig())
	r.Run(":8080")
}

func InitConfig() *Config {
	// 设置配置文件路径，根据环境切换
	var configFilePath string
	if os.Getenv("ENV") == "production" {
		configFilePath = "configs/config.yaml" // 正式环境
	} else {
		configFilePath = "configs/config-dev.yaml" // 测试环境
	}

	// 加载配置文件
	config, err := loadConfig(configFilePath)
	if err != nil {
		fmt.Println("123")
		log.Fatal("Failed to load config:", err)
	}
	return config
}

// 初始化并返回一个 Gin 路由器
func setupRouter(config *Config) *gin.Engine {
	r := gin.Default()
	var err error

	// 连接 MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.MySQL.User, config.MySQL.Password, config.MySQL.Host, config.MySQL.Port, config.MySQL.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}

	// 自动迁移
	db.AutoMigrate(&User{})

	// 连接 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password, // Redis 密码
		DB:       0,                     // 默认数据库
	})

	// 连接 Kafka
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(config.Kafka.Brokers, kafkaConfig)
	if err != nil {
		log.Fatal("Failed to connect to Kafka:", err)
	}
	defer producer.Close()

	// 路由：接受查询参数并执行 MySQL、Redis 和 Kafka 操作
	r.GET("/user", func(c *gin.Context) {
		userIDStr := c.DefaultQuery("user_id", "0")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user_id",
			})
			return
		}

		// MySQL：查询用户信息
		var user User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		// Redis：增加访问次数
		ctx := context.Background()
		visitCountKey := fmt.Sprintf("user:%d:visits", userID)
		_, err = rdb.Incr(ctx, visitCountKey).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to increment visit count in Redis",
			})
			return
		}

		// Kafka：记录用户访问事件
		event := fmt.Sprintf("User %d accessed at %s", userID, time.Now().Format(time.RFC3339))
		kafkaMsg := &sarama.ProducerMessage{
			Topic: "user-events",
			Value: sarama.StringEncoder(event),
		}
		_, _, err = producer.SendMessage(kafkaMsg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to send event to Kafka",
			})
			return
		}

		// 返回用户信息和访问次数
		visitCount, _ := rdb.Get(ctx, visitCountKey).Result()
		c.JSON(http.StatusOK, gin.H{
			"user":        user,
			"visit_count": visitCount,
		})
	})

	return r
}
