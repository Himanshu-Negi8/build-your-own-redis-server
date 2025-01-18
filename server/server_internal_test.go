package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"testing"
)

func TestSever(t *testing.T) {
	// Test server here

	serverAddr := flag.String("address", "127.0.0.1:6379", "Redis server address")
	totalConnections := flag.Int("connections", 10000, "Number of concurrent connections")
	flag.Parse()

	fmt.Printf("Testing with %d connections to %s...\n", *totalConnections, *serverAddr)

	wg := sync.WaitGroup{}
	var mu sync.Mutex
	var successCount, failureCount int

	for i := 0; i < *totalConnections; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", *serverAddr)
			if err != nil {
				mu.Lock()
				failureCount++
				mu.Unlock()
				log.Printf("Connection %d failed: %v\n", id, err)
				return
			}
			defer conn.Close()

			_, err = conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
			if err != nil {
				mu.Lock()
				failureCount++
				mu.Unlock()
				log.Printf("Connection %d write failed: %v\n", id, err)
				return
			}

			buf := make([]byte, 1024)
			_, err = conn.Read(buf)
			if err != nil {
				mu.Lock()
				failureCount++
				mu.Unlock()
				log.Printf("Connection %d read failed: %v\n", id, err)
				return
			}

			if string(buf[:7]) == "+PONG\r\n" {
				mu.Lock()
				successCount++
				mu.Unlock()
			} else {
				mu.Lock()
				failureCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("Test complete: %d successes, %d failures\n", successCount, failureCount)
}
