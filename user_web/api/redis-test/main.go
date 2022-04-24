package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
	})
	r, err := rdb.Get("18782222220").Result()
	if err != nil{
		panic(err)
	}
	fmt.Println(r)
}


