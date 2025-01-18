# Build Your Own Redis Server

This project is a simple implementation of a Redis-like server in Go. It supports basic Redis commands such as `PING`, `ECHO`, `SET`, `GET`, `CONFIG`, and `SAVE`.

## Features
- Supports basic Redis commands.
- Graceful shutdown to ensure all resources are properly cleaned up.

## Getting Started

### Prerequisites

- Understanding of Go and Redis. 
- Redis client to interact with the server (e.g., `redis-cli`).
- Go installed on your machine.

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/Himanshu-Negi8/build-your-own-redis-server.git
   cd build-your-own-redis-server
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

### Running the Server

To start the server, run:
```sh
go run server/server.go
```

The server will start listening on `0.0.0.0:6379`.

## Usage

You can interact with the server using any Redis client. Here are some example commands:

- **PING**:
  ```sh
  redis-cli -h 127.0.0.1 -p 6379 ping
  ```

- **ECHO**:
  ```sh
  redis-cli -h 127.0.0.1 -p 6379 echo "Hello, World!"
  ```

- **SET**:
  ```sh
  redis-cli -h 127.0.0.1 -p 6379 set mykey "myvalue"
  ```

- **GET**:
  ```sh
  redis-cli -h 127.0.0.1 -p 6379 get mykey
  ```

- **CONFIG**:
  ```sh
  redis-cli -h 127.0.0.1 -p 6379 config get dir
  redis-cli -h 127.0.0.1 -p 6379 config get dbfilename
  ```

- **SAVE**:
  ```sh
  redis-cli -h 127.0.0.1 -p 6379 save
  ```
I am using a docker container to interact with the server. You can use the same or use the redis-cli on your local machine. The docker container has the redis-cli installed in it.
```
docker run --rm -it redis:alpine redis-cli -h host.docker.internal -p 6379
```

## Project Structure

- `server/server.go`: Main server implementation.
- `handler/handler.go`: Command handlers for the Redis commands.
- `parser/parser.go`: RESP protocol parser.
- `types/types.go`: Custom types used in the project.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
