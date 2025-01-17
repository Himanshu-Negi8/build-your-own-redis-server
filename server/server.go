package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/Himanshu-Negi8/build-your-own-redis-server/handler"
	"github.com/Himanshu-Negi8/build-your-own-redis-server/parser"
	"github.com/Himanshu-Negi8/build-your-own-redis-server/types"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	cache := make(map[string]types.CustomValue)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("failed to serve requests")
			continue
		}

		conn.SetReadDeadline(time.Now().Add(time.Second * 10)) // Timeout after 10 seconds
		go handleConnectionRequest(conn, cache)
	}

}

func handleConnectionRequest(conn net.Conn, cache map[string]types.CustomValue) {
	defer conn.Close()
	buf := make([]byte, 2048)

	// The reason for this infinite loop is to handle all the incoming commands within the same connection.
	// Obviously since it's socket connection, it will be closed by the client after the client is done with the commands.
	for {
		// Read the client's input
		_, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// Client closed the connection.
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
