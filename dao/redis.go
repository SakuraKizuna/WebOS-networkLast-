package dao

import (
	"github.com/go-redis/redis"
	"math/rand"
	"time"
)

var (
	Redis *redis.Client
)

// TODO dao.Redis.Expire(token,1200*time.Second)重置key的时间

func InitRedis() (err error) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "020804", // no password set
		//Password:"123456",
		DB:       0,        // use default DB
	})
	_, err = Redis.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
