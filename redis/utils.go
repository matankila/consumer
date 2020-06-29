package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

func NewBD(c *redis.Client, bdName string, body string) {
	err := c.Set("bd-"+bdName, "success", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Printf("bd created: %s, from the settings: %s\n", bdName, body)
}
