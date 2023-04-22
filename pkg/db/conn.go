// Copyright 2023 N4-Networks.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbCfg struct {
	serverAddr string
	name       string
	userName   string
	passwd     string
	timeout    int // in minute
}

var cfg dbCfg

func readConfigFromEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Println("Error in loading .env file")
		return err
	}

	// DB config params
	if env, ok := os.LookupEnv("DB_ADDR"); ok {
		cfg.serverAddr = env
	} else {
		log.Println("DB address is not set")
		return errors.New("DB address is not set")
	}

	if env, ok := os.LookupEnv("DB_USER"); ok {
		cfg.userName = env
	} else {
		log.Println("DB User name is not set")
		return errors.New("DB_USER is not set")
	}

	if env, ok := os.LookupEnv("DB_PASSWD"); ok {
		cfg.passwd = env
	} else {
		log.Println("DB User password is not set")
		return errors.New("DB_PASSWD is not set")
	}

	if env, ok := os.LookupEnv("DB_NAME"); ok {
		cfg.name = env
	} else {
		log.Println("DB Name is not set")
		return errors.New("DB_NAME is not set")
	}

	if env, ok := os.LookupEnv("DB_CONN_TIMEOUT"); ok {
		x, _ := strconv.ParseInt(env, 10, 0)
		cfg.timeout = int(x)
	} else {
		cfg.timeout = 3
		log.Println("DB Connection Timeout is not set, default 3mins")
	}
	log.Printf("DB Config params: %+v\n", cfg)
	return nil

}

func Connect() (*mongo.Client, error) {
	if err := readConfigFromEnv(); err != nil {
		return nil, err
	}
	cred := options.Credential{Username: cfg.userName, Password: cfg.passwd}
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + cfg.serverAddr).SetAuth(cred))
	if err != nil {
		return nil, err
	}
	timeout := time.Duration(cfg.timeout) * time.Minute
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	if err = client.Connect(ctx); err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	return client, err
}

/*
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
*/

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
