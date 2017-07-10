package main

import (
	"fmt"
	"time"

	"github.com/tevino/go-zookeeper/zk"
	"flag"
	"strings"
)

func main() {
	zookeeperServerArg := flag.String("zookeeperServers", "127.0.0.1:2181", "list zookeeper broker.")
	flag.Parse()
	zookeeperServers := strings.Split(*zookeeperServerArg, ",")

	action := flag.Args()[0]
	director := flag.Args()[1]

	fmt.Println(zookeeperServers)

	connect, _, err := zk.Connect(zookeeperServers, time.Second * 4)
	if err != nil {
		panic(err)
	}
	defer connect.Close()

	switch action {
	case "create":
		createDir(connect, director)
	case "delete":
		deleteDir(connect, director)
	default:
		panic(fmt.Sprint("action %s not supported", action))
	}
}

func createDir(connect *zk.Conn, director string) {
	_, err := connect.Create(director, nil, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func deleteDir(connect *zk.Conn, director string) {
	err := connect.Delete(director, -1)
	if err != nil {
		fmt.Println(err)
		return
	}
}