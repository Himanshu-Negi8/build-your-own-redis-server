package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/Himanshu-Negi8/build-your-own-redis-server/handler"
	"github.com/Himanshu-Negi8/build-your-own-redis-server/parser"
	"github.com/Himanshu-Negi8/build-your-own-redis-server/types"
)

const workers = 10

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	cache := make(map[string]types.CustomValue)
	connCh := make(chan net.Conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 1; i <= workers; i++ {
		go worker(ctx, connCh, cache)
	}

	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Println("failed to serve requests")
			continue
		}

		connCh <- conn

	}

}

func worker(ctx context.Context, connCh chan net.Conn, cache map[string]types.CustomValue) {
	for {
		select {
		case <-ctx.Done():
			return
		case conn := <-connCh:
			fmt.Println("worker launching for connection")
			// calling handleConnectionRequest function as a goroutine to make sure that the worker is not blocked
			go handleConnectionRequest(ctx, conn, cache)
		}
	}
}

func handleConnectionRequest(ctx context.Context, conn net.Conn, cache map[string]types.CustomValue) {
	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			buf := make([]byte, 8196)
			// Read the client's input
			_, err := conn.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					// Client closed the connection
					fmt.Println("Client disconnected")
					return
				}
				fmt.Println("Error reading from connection:", err)
				return
			}

			// Parse the RESP command
			r := bytes.NewReader(buf)
			result, respType, err := parser.Parse(r)
			if err != nil {
				fmt.Println("Error parsing RESP:", err)
				conn.Write([]byte("-ERR invalid command\r\n"))
				continue
			}

			// Process the command based on its type
			response := handler.HandleCommands(result, respType, cache)
			fmt.Println(string(response))
			conn.Write(response)
		}
	}
}
