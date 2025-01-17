# build-your-own-redis-server

## Build Your Own Redis Server
A lightweight implementation of a Redis-like server from scratch, focusing on the RESP protocol and a custom interface for core Redis commands. Perfect for learning, experimenting, and understanding the internals of Redis.

### Current Pain Point Or Area to Explore
I attempted to implement the Go concurrency pattern using a fan-out approach to handle multiple requests with a configurable number of workers. However, it failed to accept or establish connection requests beyond a certain number of workers if you don't call the handleConncetionRequest concurrently as it blocks the worker. Each worker was supposed to maintain an open connection, as it’s a socket connection, and the client could submit multiple commands during one connection. Afterward, I reverted to directly calling the handler without using a worker pool. Moving forward, I might try implementing an event-loop mechanism, similar to how Redis handles it.
