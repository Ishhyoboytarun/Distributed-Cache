package main

import (
    "fmt"
    "log"
    "sync"

    "github.com/hashicorp/consul/api"
)

type DistributedCache struct {
    client *api.Client
    mutex  sync.Mutex
}

func NewDistributedCache() (*DistributedCache, error) {
    client, err := api.NewClient(&api.Config{
        Address: "localhost:8500",
    })
    if err != nil {
        return nil, err
    }

    return &DistributedCache{
        client: client,
        mutex:  sync.Mutex{},
    }, nil
}

func (dc *DistributedCache) Get(key string) (string, error) {
    dc.mutex.Lock()
    defer dc.mutex.Unlock()

    pair, _, err := dc.client.KV().Get(key, nil)
    if err != nil {
        return "", err
    }
    if pair == nil {
        return "", fmt.Errorf("key %s does not exist in cache", key)
    }

    return string(pair.Value), nil
}

func (dc *DistributedCache) Set(key string, value string) error {
    dc.mutex.Lock()
    defer dc.mutex.Unlock()

    pair := &api.KVPair{
        Key:   key,
        Value: []byte(value),
    }
    _, err := dc.client.KV().Put(pair, nil)
    return err
}

func main() {
    dc, err := NewDistributedCache()
    if err != nil {
        log.Fatal(err)
    }

    err = dc.Set("foo", "bar")
    if err != nil {
        log.Fatal(err)
    }

    val, err := dc.Get("foo")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(val)
}
