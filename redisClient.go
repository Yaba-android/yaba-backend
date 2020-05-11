package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

/********************************
*
* 		DATABASE REDIS
*
********************************/

func redisSetNewBook(client *redis.Client, book *Book) error {
	m := make(map[string]interface{})
	id := uuid.NewV5(uuid.NamespaceOID, book.Title)

	book.RemoteId = id.String()
	bookJSON, err := json.Marshal(book)
	if err != nil {
		fmt.Println(err)
		return err
	}
	m[book.RemoteId] = bookJSON
	err = client.HMSet(RedisTableBooks, m).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func redisGetAllBooks(client *redis.Client) []Book {
	var bookObj Book
	var booksSlice []Book

	booksString, err := client.HGetAll(RedisTableBooks).Result()
	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}
	for _, bookStr := range booksString {
		bookBytes := []byte(bookStr)
		err = json.Unmarshal(bookBytes, &bookObj)
		booksSlice = append(booksSlice, bookObj)
	}
	//printBookSlice(booksSlice)
	return booksSlice
}

func redisStartClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv(RedisHostAddr),
		Password: os.Getenv(RedisDbPassword),
		DB:       0,
	})
	return client
}

func redisIsClientConnected(client *redis.Client) error {
	pong, err := client.Ping().Result()

	if err != nil {
		fmt.Println("Redis error")
	} else if pong == RedisPong {
		fmt.Println("Redis connected")
	}
	return err
}
