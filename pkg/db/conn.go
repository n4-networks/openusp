package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(addr string, user string, passwd string, timeout time.Duration) (*mongo.Client, error) {
	cred := options.Credential{Username: user, Password: passwd}
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + addr).SetAuth(cred))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	if err = client.Connect(ctx); err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	return client, err
}

func ConnectCache(addr string, timeout time.Duration) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Println("Error in connecting to redis", addr)
		return nil, err
	}
	fmt.Println(pong)
	return client, nil

}
