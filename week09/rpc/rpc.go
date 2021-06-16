package main

import (
	"log"
	"net/rpc"
)

func main() {

	serverAddress := "localhost"
	_, err := rpc.DialHTTP("tcp", serverAddress + ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
}
