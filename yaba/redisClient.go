package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

/********************************
*
* 		DATABASE REDIS
*
********************************/

// Book struct
type Book struct {
	Name    string
	Author  string
	DatePub string
	Path    string
}

func redisSetNewBook(client *redis.Client, book *Book) error {
	bookJSON, err := json.Marshal(book)
	if err != nil {
		fmt.Println(err)
		return err
	}
	id := uuid.NewV5(uuid.NamespaceOID, book.Name)
	err = client.HMSet("books", id, bookJSON).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func redisGetAllBooks(client *redis.Client) error {
	books, err := client.HGetAll("books").Result()
	if err != nil {
		fmt.Println("error", err)
		return err
	}
	printMap(books)
	return err
}

func redisStartClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return client
}

func redisIsClientConnected(client *redis.Client) error {
	pong, err := client.Ping().Result()

	if err != nil {
		fmt.Println("error")
	} else if pong == "PONG" {
		fmt.Println("Connected")
	}
	return err
}
