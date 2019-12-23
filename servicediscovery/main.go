package main

import (
	"flag"
	"fmt"

	"github.com/zhangkun/etcd-example/servicediscovery/discovery"
)

func main() {
	var role = flag.String("role", "", "master | worker")
	flag.Parse()
	endpoints := []string{"http://10.151.3.8:2379"}
	if *role == "master" {
		master := discovery.NewMaster(endpoints)
		master.WatchWorkers()
	} else if *role == "worker" {
		worker := discovery.NewWorker("localhost", "10.151.3.8", endpoints)
		worker.HeartBeat()
	} else {
		fmt.Println("example -h for usage")
	}
}
