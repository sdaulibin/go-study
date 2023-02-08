package models

import (
	"context"
	"ginchat/utils"
	"time"
)

/*
*
设置在线用户到redis缓存
*
*/
func SetUserOnlineInfo(key string, data []byte, timeTTL time.Duration) {
	ctx := context.Background()
	utils.RedisClient.Set(ctx, key, data, timeTTL)
}
