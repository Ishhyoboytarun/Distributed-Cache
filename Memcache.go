package main

import (
    "fmt"
    "log"
    "sync"

    "github.com/bradfitz/gomemcache/memcache"
)

type DistributedCache struct {
    client *memcache.Client
    mutex  sync.Mutex
}

func NewDistributedCache() *DistributedCache {
    client := memcache.New("localhost:11211")

    return &DistributedCache{
        client: client,
        mutex:  sync.Mutex{},
    }
}

func (dc *DistributedCache) Get(key string) (string, error) {
    dc.mutex.Lock()
    defer dc.mutex.Unlock()

    item, err := dc.client.Get(key)
    if err != nil {
        return "", fmt.Errorf("key %s does not exist in cache", key)
    }

    return string(item.Value), nil
}

func (dc *DistributedCache) Set(key string, value string) error {
    dc.mutex.Lock()
    defer dc.mutex.Unlock()

    item := memcache.Item{
        Key:   key,
        Value: []byte(value),
    }
    return dc.client.Set(&item)
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
