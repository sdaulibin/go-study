package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

const (
	PUBLISH_KEY = "websocket"
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app inited .......")
}

func InitMysql() {
	sqlLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{
		Logger: sqlLogger,
	})
	fmt.Println("mysql inited .......")
}

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	//pong, err := REDIS_CLIENT.Ping().Result()
	// if err != nil {
	// 	fmt.Println("init redis client err：", err)
	// 	return
	// }
	fmt.Println("init redis client success ......")
}

func Publish(ctx context.Context, channel string, msg string) error {
	fmt.Println("Publish>>>>>>", msg)
	err := RedisClient.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := RedisClient.Subscribe(ctx, channel)
	fmt.Println("Subscribe>>>>>>", sub)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Subscribe msg >>>>>>", msg)
	return msg.Payload, err
}
