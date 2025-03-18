package main

import (
	"github.com/Himanshu-Negi8/build-your-own-redis-server/server"
	"log"
)

func main() {
	s := server.NewServer("127.0.0.1", 6379, 2000)
	err := s.RunAsyncServer()
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
