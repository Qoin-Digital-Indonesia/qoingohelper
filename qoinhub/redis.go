package qoingohelper

import (
	"context"
	"encoding/json"
	"log"

	"time"

	redis "github.com/go-redis/redis/v8"
)

var redisPoolClient *redis.Client

func GetRedisClient(redisHost, redisPort, redisPassword string, redisDb int) (*redis.Client, error) {

	ctx := context.Background()
	redisPoolClient = redis.NewClient(&redis.Options{
		Addr:       redisHost + ":" + redisPort,
		Password:   redisPassword,
		DB:         redisDb,
		MaxRetries: 3,
		PoolSize:   50,
		MaxConnAge: 5 * time.Minute,
	})

	res, err := redisPoolClient.Ping(ctx).Result()
	if err != nil {
		LoggerErrorHub(err)
		return nil, err
	}

	log.Println("open redis pool connection successfully")
	log.Println(res)

	return redisPoolClient, nil
}

func StoreRedis(id string, data interface{}, duration time.Duration) (err error) {

	client := redisPoolClient
	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		LoggerErrorHub("error redis ping : " + err.Error())
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = client.Set(ctx, id, string(jsonData), duration).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetRedis(id string) (result string, err error) {

	client := redisPoolClient
	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		LoggerErrorHub("error redis ping : " + err.Error())
		return
	}

	ctx := context.Background()
	getRedis := client.Get(ctx, id)
	if getRedis == nil {
		return
	}

	if err = getRedis.Err(); err != nil {
		return
	}

	return getRedis.Result()
}

func DeleteRedis(id string) (err error) {

	client := redisPoolClient
	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		LoggerErrorHub("error redis ping : " + err.Error())
		return
	}

	delete := client.Del(context.Background(), id)
	if delete == nil {
		return
	}

	if err = delete.Err(); err != nil {
		return
	}

	return
}
