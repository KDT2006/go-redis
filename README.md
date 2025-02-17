# Go Redis Server

A minimal Redis-compatible server implementation in Go, built from scratch with only one external dependency ([tidwall/resp](https://github.com/tidwall/resp)) for RESP serialization. This project was created to deepen my understanding of networking, concurrency, and in-memory data storage.

## Features

- **Redis Protocol Support**: Implements the RESP (REdis Serialization Protocol) using the [tidwall/resp](https://github.com/tidwall/resp) library.
- **Basic Command Support**: Supports a subset of Redis commands, including:
  - `GET`: Retrieve the value of a key.
  - `SET`: Set the value of a key.
  - `CLIENT`: Client connection management.
  - `HELLO`: Handshake with the server.
- **Concurrency Handling**: Utilizes Go's goroutines and channels to handle multiple client connections efficiently.
- **Tested with Redis Client**: Verified compatibility with the official Redis CLI.

## Getting Started

### Prerequisites

- Go 1.20 or higher
- Basic understanding of Redis and its commands

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/KDT2006/go-redis.git
   cd go-redis
   ```

2. Build the server:

   ```
   go build -o goredis
   ```

3. Run the server:
   ```
   ./goredis
   ```

## Why This Project?

This project was built to:

- **Learn Go Concurrency**: Deepen my understanding of Go's concurrency model, including goroutines and channels, by handling multiple client connections efficiently.
- **Understand Redis Internals**: Gain hands-on experience with how Redis works under the hood, including its protocol (RESP) and in-memory data storage.
- **Explore Networking**: Learn about TCP servers, client-server communication, and protocol implementation in a real-world context.
- **Showcase Problem-Solving Skills**: Demonstrate my ability to build a functional, minimalistic system from scratch with minimal dependencies.

## Future Improvements

Here are some potential enhancements for this project:

- **Support More Redis Commands**: Add support for additional Redis commands like `DEL`, `EXISTS`, `INCR`, `EXPIRE`.
- **Benchmarking**: Measure and optimize the server's performance under high load.
