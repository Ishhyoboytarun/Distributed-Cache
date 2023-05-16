package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-redis/redis"
)

type DistributedCache struct {
	client *redis.Client
	mutex  sync.Mutex
}

func NewDistributedCache() *DistributedCache {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &DistributedCache{
		client: client,
		mutex:  sync.Mutex{},
	}
}

func (dc *DistributedCache) Get(key string) (string, error) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	val, err := dc.client.Get(key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s does not exist in cache", key)
	} else if err != nil {
		return "", err
	}

	return val, nil
}

func (dc *DistributedCache) Set(key string, value string) error {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	return dc.client.Set(key, value, 0).Err()
}

func main() {
	dc := NewDistributedCache()

	err := dc.Set("foo", "bar")
	if err != nil {
		log.Fatal(err)
	}

	val, err := dc.Get("foo")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(val)
}
