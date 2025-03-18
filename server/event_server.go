package server

import (
	"bytes"
	"errors"
	"log"
	"net"
	"syscall"
	"time"

	"github.com/Himanshu-Negi8/build-your-own-redis-server/handler"
	"github.com/Himanshu-Negi8/build-your-own-redis-server/iomultiplexer"
	"github.com/Himanshu-Negi8/build-your-own-redis-server/parser"
	"github.com/Himanshu-Negi8/build-your-own-redis-server/types"
)

type server struct {
	serverFD    int    // file descriptor of the server
	host        string // ip address only
	port        int    // port number only
	maxClients  int    // maximum number of clients that can connect to the server
	multiplexer iomultiplexer.IOMultiplexer
	cache       map[string]types.CustomValue
}

func NewServer(host string, port, maxClients int) *server {
	return &server{
		host:       host,
		port:       port,
		maxClients: maxClients,
		cache:      make(map[string]types.CustomValue),
	}
}

func (s *server) RunAsyncServer() error {
	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		// handle the error
		return err
	}

	s.serverFD = serverFD
	// Close the socket on exit if an error occurs
	defer func() {
		if err != nil {
			if err := syscall.Close(serverFD); err != nil {
				log.Print("failed to close server socket", "error", err)
			}
		}
	}()

	if err := syscall.SetsockoptInt(serverFD, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return err
	}

	err = syscall.SetNonblock(serverFD, true)
	if err != nil {
		// handle the error
		return err
	}

	ip4 := net.ParseIP(s.host)
	syscall.Bind(serverFD, &syscall.SockaddrInet4{Port: s.port, Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]}})

	defer func(fd int) {
		err := syscall.Close(fd)
		if err != nil {
			log.Println("failed to close server socket", "error", err)
		}
	}(s.serverFD)

	err = syscall.Listen(serverFD, s.maxClients)
	if err != nil {
		return err
	}

	s.multiplexer, err = iomultiplexer.New(s.maxClients)
	if err != nil {
		// handle the error
		return err
	}

	if err := s.multiplexer.Subscribe(iomultiplexer.Event{
		Fd: s.serverFD,
		Op: iomultiplexer.OpRead,
	}); err != nil {
		return err
	}

	err = s.eventLoop()
	return err
}

func (s *server) eventLoop() error {
	for {
		events, err := s.multiplexer.Poll(time.Millisecond * 100)
		if err != nil {
			if errors.Is(err, syscall.EINTR) {
				continue
			}
			return err
		}

		for _, event := range events {
			if event.Fd == s.serverFD {
				if err := s.acceptClientConnection(); err != nil {
					log.Println("failed to accept client connection", "error", err)
				}
			} else {
				err := s.handleClientEvent(event)
				if err != nil {
					log.Println("failed to handle client event", "error", err)
				}

			}
		}

	}

}

// acceptClientConnection accepts a new client connection and subscribes to read events on the connection.
func (s *server) acceptClientConnection() error {
	fd, _, err := syscall.Accept(s.serverFD)
	if err != nil {
		return err
	}

	if err := syscall.SetNonblock(fd, true); err != nil {
		return err
	}

	return s.multiplexer.Subscribe(iomultiplexer.Event{
		Fd: fd,
		Op: iomultiplexer.OpRead,
	})
}

// handleClientEvent reads commands from the client connection and responds to the client. It also handles disconnections.
func (s *server) handleClientEvent(event iomultiplexer.Event) error {
	// Read from the file descriptor
	buf := make([]byte, 4096)
	n, err := syscall.Read(event.Fd, buf)
	if err != nil {
		if errors.Is(err, syscall.EAGAIN) {
			// No data available, return nil
			return nil
		}
		// Handle other read errors
		return err
	}

	r := bytes.NewReader(buf[:n])
	result, respType, err := parser.Parse(r)
	if err != nil {
		log.Printf("Error parsing RESP: %v", err)
		response := []byte("-ERR invalid command\r\n")
		syscall.Write(event.Fd, response)
		return err
	}

	// Process the command based on its type
	response := handler.HandleCommands(result, respType, s.cache)
	_, err = syscall.Write(event.Fd, response)
	if err != nil {
		return err
	}

	return nil
}
