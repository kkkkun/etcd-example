package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/etcd-io/etcd/client"
)

func NewClient(endpoints []string) client.KeysAPI {
	cfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: 10 * time.Second,
	}
	etcdclient, err := client.New(cfg)
	if err != nil {
		log.Fatal("Error:cannot connec to etcd:", err)
	}
	client := client.NewKeysAPI(etcdclient)
	return client
}

func main() {
	endpoints := []string{"http://10.151.3.8:2379"}
	producer := NewClient(endpoints)
	consumer := NewClient(endpoints)

	watcherOptions := &client.WatcherOptions{
		Recursive:  false,
		AfterIndex: 0,
	}
	etcdclt := consumer.Watcher("say", watcherOptions)

	n := 1
	go func() {
		for {
			producer.Set(context.Background(), "say", strconv.Itoa(n), &client.SetOptions{})
			n++
			time.Sleep(2 * time.Second)
		}
	}()

	for {
		resp, err := etcdclt.Next(context.TODO())
		if err != nil {
			fmt.Println("Error,consumer err")
		}
		if resp.Node.Dir {
			continue
		}
		fmt.Printf("[%s] %s %s\n", resp.Action, resp.Node.Key, resp.Node.Value)
	}
}
