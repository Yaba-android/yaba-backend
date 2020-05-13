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

func redisSetNewAuthor(client *redis.Client, author *Author) *Author {
	m := make(map[string]interface{})
	id := uuid.NewV5(uuid.NamespaceOID, author.Name)

	author.RemoteId = id.String()
	authorJSON, err := json.Marshal(author)
	if err != nil {
		fmt.Println("setNewAuthor Marshal error: ", err)
		return nil
	}
	m[author.RemoteId] = authorJSON
	err = client.HMSet(RedisTableAuthors, m).Err()
	if err != nil {
		fmt.Println("setNewAuthor error: ", err)
		return nil
	}
	return author
}

func redisGetAllAuthors(client *redis.Client) []Author {
	var authorObj Author
	var authorsSlice []Author

	authorsString, err := client.HGetAll(RedisTableAuthors).Result()
	if err != nil {
		fmt.Println("getAllAuthors error: ", err)
		return nil
	}
	for _, authorStr := range authorsString {
		authorBytes := []byte(authorStr)
		err = json.Unmarshal(authorBytes, &authorObj)
		if err != nil {
			fmt.Println("getAllAuthors Unmarshal error: ", err)
			return nil
		}
		authorsSlice = append(authorsSlice, authorObj)
	}
	return authorsSlice
}

func redisGetAuthorById(client *redis.Client, remoteId string) *Author {
	var authorObj Author

	authorStr, err := client.HGet(RedisTableAuthors, remoteId).Result()
	if err != nil {
		fmt.Println("getAuthor error: ", err)
		return nil
	}
	authorBytes := []byte(authorStr)
	err = json.Unmarshal(authorBytes, &authorObj)
	if err != nil {
		fmt.Println("getAuthor Unmarshal error: ", err)
		return nil
	}
	return &authorObj
}

func redisSetNewBook(client *redis.Client, book *Book) *Book {
	m := make(map[string]interface{})
	id := uuid.NewV5(uuid.NamespaceOID, book.Title)

	book.RemoteId = id.String()
	bookJSON, err := json.Marshal(book)
	if err != nil {
		fmt.Println("setNewBook Marshal error: ", err)
		return nil
	}
	m[book.RemoteId] = bookJSON
	err = client.HMSet(RedisTableBooks, m).Err()
	if err != nil {
		fmt.Println("setNewBook error: ", err)
		return nil
	}
	return book
}

func redisGetAllBooks(client *redis.Client) []Book {
	var bookObj Book
	var booksSlice []Book

	booksString, err := client.HGetAll(RedisTableBooks).Result()
	if err != nil {
		fmt.Println("getAllBooks error: ", err)
		return nil
	}
	for _, bookStr := range booksString {
		bookBytes := []byte(bookStr)
		err = json.Unmarshal(bookBytes, &bookObj)
		if err != nil {
			fmt.Println("getAllBooks Unmarshal error: ", err)
			return nil
		}
		booksSlice = append(booksSlice, bookObj)
	}
	return booksSlice
}

func redisGetBookById(client *redis.Client, remoteId string) *Book {
	var bookObj Book

	bookStr, err := client.HGet(RedisTableBooks, remoteId).Result()
	if err != nil {
		fmt.Println("getBook error: ", err)
		return nil
	}
	bookBytes := []byte(bookStr)
	err = json.Unmarshal(bookBytes, &bookObj)
	if err != nil {
		fmt.Println("getBook Unmarshal error: ", err)
		return nil
	}
	return &bookObj
}

func redisSetNewGenre(client *redis.Client, genre *Genre) *Genre {
	m := make(map[string]interface{})
	id := uuid.NewV5(uuid.NamespaceOID, genre.Name)

	genre.RemoteId = id.String()
	genreJSON, err := json.Marshal(genre)
	if err != nil {
		fmt.Println("setNewGenre Marshal error: ", err)
		return nil
	}
	m[genre.RemoteId] = genreJSON
	err = client.HMSet(RedisTableGenres, m).Err()
	if err != nil {
		fmt.Println("setNewGenre error: ", err)
		return nil
	}
	return genre
}

func redisGetAllGenres(client *redis.Client) []Genre {
	var genreObj Genre
	var genresSlice []Genre

	genresString, err := client.HGetAll(RedisTableGenres).Result()
	if err != nil {
		fmt.Println("getAllGenres error: ", err)
		return nil
	}
	for _, genreStr := range genresString {
		bookBytes := []byte(genreStr)
		err = json.Unmarshal(bookBytes, &genreObj)
		if err != nil {
			fmt.Println("getAllGenres Unmarshal error: ", err)
			return nil
		}
		genresSlice = append(genresSlice, genreObj)
	}
	return genresSlice
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
