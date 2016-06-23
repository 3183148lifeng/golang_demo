package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var rp *redis.Pool

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			fmt.Println("pwd: ", password)
			if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

var (
	pool          *redis.Pool
	redisServer   = flag.String("redisServer", ":6379", "")
	redisPassword = flag.String("redisPassword", "12345678", "")
	key           = flag.String("key", "name", "")
	value         = flag.String("value", "fli", "")
)

func main() {

	flag.Parse()
	fmt.Println("key: ", *key)
	fmt.Println("value: ", *value)

	pool = newPool(*redisServer, *redisPassword)

	conn := pool.Get()

	_, err := conn.Do("SET", *key, *value)
	if err != nil {
		fmt.Println(err)
	}
	name, _ := redis.String(conn.Do("GET", *key))
	fmt.Println("name: ", name)
	// defer conn.Close()
}
