package main

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"

    "go.etcd.io/etcd/clientv3"
)

type DistributedCache struct {
    client *clientv3.Client
    mutex  sync.Mutex
}

func NewDistributedCache() (*DistributedCache, error) {
    client, err := clientv3.New(clientv3.Config{
        Endpoints:   []string{"localhost:2379"},
        DialTimeout: 5 * time.Second,
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

    resp, err := dc.client.Get(context.Background(), key)
    if err != nil {
        return "", err
    }
    if len(resp.Kvs) == 0 {
        return "", fmt.Errorf("key %s does not exist in cache", key)
    }

    return string(resp.Kvs[0].Value), nil
}

func (dc *DistributedCache) Set(key string, value string) error {
    dc.mutex.Lock()
    defer dc.mutex.Unlock()

    _, err := dc.client.Put(context.Background(), key, value)
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
